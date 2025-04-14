package main

import (
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/go-attestation/attest"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	lib "github.com/pv204-security-technologies-project/internal"
)

type MessageData struct {
	Content string `json:"content"`
}

// shared object for session data
// TODO I am not sure if the data are encrypted in the cookie itself, or if they are stored locally in some hashmap
var store *sessions.CookieStore

var db *pgx.Conn

func handleRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var combinedRegisterData lib.CombinedRegisterData
	if err := json.Unmarshal(body, &combinedRegisterData); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	secret, secret_bytes, err := generateChallengeFromRegisterData(combinedRegisterData.Data, combinedRegisterData.PreData)

	if err != nil {
		log.Printf("Error generating challenge: %v", err)
		errMsg := fmt.Sprintf("Error generating a challenge: %v", err)
		// TODO this should probably be a different code
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// set a session (cookie) with a secret
	session, _ := store.Get(r, lib.CookieName)
	session.Values["secret"] = secret
	session.Values["ek"] = combinedRegisterData.Data.PublicKey
	session.Values["username"] = combinedRegisterData.Username
	session.Save(r, w)
	w.Header().Set("Content-Type", "application/json")
	w.Write(secret_bytes)
}

func handleRegistrationResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var clientRegResponse lib.ClientResponseData
	if err := json.Unmarshal(body, &clientRegResponse); err != nil {
		log.Fatal(err) // this should return http code instead of exit imo
	}

	session, _ := store.Get(r, lib.CookieName)

	// check if the session is set and that the secret is not none
	if len(session.Values) == 0 || len(session.Values["secret"].([]byte)) == 0 {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// get the secret from the session
	secret := session.Values["secret"].([]byte)

	if !bytes.Equal(secret, clientRegResponse.Secret) {
		// fail
		http.Error(w, "Verification failed", http.StatusUnauthorized)
		return
	}

	var name string
	err = db.QueryRow(context.Background(), "SELECT name from subject where ek = $1", session.Values["ek"]).Scan(&name)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if name != "" {
		session.Values["username"] = name
	} else {
		_, err = db.Exec(context.Background(), `INSERT INTO subject (name, ek)
      VALUES ($1, $2)
      RETURNING id`, session.Values["username"], session.Values["ek"])

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Print(err)
			return
		}
	}

	session.Values["Authenticated"] = true
	session.Save(r, w)

	w.Write([]byte(fmt.Sprintf("Welcome back %s", name)))
}

func pemToPublicKey(pubPEM string) (crypto.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

func generateChallengeFromRegisterData(data lib.RegisterData, preData lib.PreRegisterData) ([]byte, []byte, error) {

	cryptoPubKey, err := pemToPublicKey(data.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing publicKey: %v", err)
	}

	var params attest.ActivationParameters

	switch preData.PublicKeyType {
	case lib.RSA:
		params = attest.ActivationParameters{
			TPMVersion: data.Version,
			EK:         cryptoPubKey,
			AK:         data.AttestParams,
		}
	case lib.ECDH:
		params = attest.ActivationParameters{
			TPMVersion: data.Version,
			EK:         cryptoPubKey,
			AK:         data.AttestParams,
		}
	case lib.ECDSA:
		params = attest.ActivationParameters{
			TPMVersion: data.Version,
			EK:         cryptoPubKey,
			AK:         data.AttestParams,
		}
	default:
		log.Fatal("Unsupported key type")
	}
	secret, encryptedCredentials, err := params.Generate()
	lib.Handle_error(err)

	// "send" the response to the client
	secret_bytes, err := json.Marshal(encryptedCredentials)
	lib.Handle_error(err)

	// this is the secret that the client needs to uncover
	return secret, secret_bytes, nil
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, lib.CookieName)

	// check authentication
	if session.Values["Authenticated"] == nil || session.Values["Authenticated"] == false {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	var msgData MessageData
	err := json.NewDecoder(r.Body).Decode(&msgData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Message received: %s\n", msgData.Content)

	//logging the message with time
	f, err := os.OpenFile("board_of_messages.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(time.Now().Format(time.RFC850) + ": " + msgData.Content + "\n"); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Message received successfully")
}

func downloadMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, lib.CookieName)

	if session.Values["Authenticated"] == nil || session.Values["Authenticated"] == false {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	data, err := os.ReadFile("board_of_messages.txt")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Could load the messages")
		return
	}

	msgData := lib.MessageData{
		Content: string(data),
	}

	jsonData, err := json.Marshal(msgData)
	if err != nil {
		log.Fatal("Error marshalling message data", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// init function is a special function, just like main is a bit special
// the functin is called before main to initialize objects
func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	// I am not sure what this piece of code does really
	// and what is the significance of the keys
	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
}

func setupDatabase() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_CONNECTION_STRING"))
	if err != nil {
		fmt.Print(err)
		log.Fatal("Unable to connect to database.")
	}

	db = conn

	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS subject (
      id    SERIAL NOT NULL PRIMARY KEY,
      name  TEXT   NOT NULL UNIQUE,
      ek    TEXT   NOT NULL UNIQUE);`)

	if err != nil {
		log.Fatal(err)
	}
}

func setupEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, watch out for errors")
	}
}

func main() {
	// TODO these paths should be set as constants in a library and shared
	// between the client and the server

  setupEnv()
  setupDatabase()

	http.HandleFunc("/registration", handleRegistration)
	http.HandleFunc("/registration/finish", handleRegistrationResponse)
	http.HandleFunc("/message", handleMessage)
	http.HandleFunc("/download", downloadMessages)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServeTLS(":8080", "tls/server.crt", "tls/server.key", nil))
}
