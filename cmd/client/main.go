package main

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/google/go-attestation/attest"
	lib "github.com/pv204-security-technologies-project/internal"
)

func publicKeyToPEM(pub crypto.PublicKey) (string, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})
	return string(pubPEM), nil
}

// client_generate_ek_ak generates and returns the EK and AK from the TPM.
func client_generate_ek_ak() (*attest.AK, *lib.CombinedRegisterData, error) {
	// Open the TPM session.
	tpm, err := attest.OpenTPM(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening TPM: %v", err)
	}
	defer tpm.Close()

	// Get the TPM version.
	version := tpm.Version()

	// TPMs are provisioned with a set of EKs by the manufacturer.
	eks, err := tpm.EKs()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting EKs: %v", err)
	}
	ek := eks[0] // Select the first EK.

	var publicKeyType lib.PublicKeyType
	switch ek.Public.(type) {
	case *ecdsa.PublicKey:
		publicKeyType = lib.ECDSA
	case *rsa.PublicKey:
		publicKeyType = lib.RSA
	case *ecdh.PublicKey:
		publicKeyType = lib.ECDH
	default:
		return nil, nil, fmt.Errorf("unsupported public key type")
	}

	// Create an AK.
	ak, err := tpm.NewAK(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating AK: %v", err)
	}
	stringPEM, err := publicKeyToPEM(ek.Public)
	lib.Handle_error(err)
	attestParams := ak.AttestationParameters()
	// Serialize the data to be sent to the server.
	preRegisterData := lib.PreRegisterData{PublicKeyType: publicKeyType}
	data := lib.RegisterData{
		Version:      version,
		PublicKey:    stringPEM,
		AttestParams: attestParams,
	}
	combinedRegisterData := lib.CombinedRegisterData{
		PreData:  preRegisterData,
		Data:     data,
		Username: "",
	}

	return ak, &combinedRegisterData, nil
}

func client_proof(challenge_bytes []byte, ak *attest.AK) []byte {

	// open TPM session
	tpm, err := attest.OpenTPM(nil)
	lib.Handle_error(err)

	// close the tpm after the end of the process
	defer tpm.Close()

	akBytes, err := ak.Marshal()
	lib.Handle_error(err)

	// load AK we generated in the previous step on the client
	pak, err := tpm.LoadAK(akBytes)
	lib.Handle_error(err)

	// read data sent from the server
	var encryptedCredentials attest.EncryptedCredential
	if err := json.Unmarshal(challenge_bytes, &encryptedCredentials); err != nil {
		log.Fatal(err)
	}

	// decrypt the secret
	secret, err := pak.ActivateCredential(tpm, encryptedCredentials)
	if err != nil {
		log.Fatal("Cred activation failed: ", err)
	}

	// send the secret back to the server as a proof that  we can
	return secret
}

// authenticate sends the registration data to the server.
func authenticate(urlprefix string, client *http.Client, combinedRegisterData lib.CombinedRegisterData, ak *attest.AK, username string) {

	url := urlprefix + "/registration"
	combinedRegisterData.Username = username

	combinedRegisterDataJSON, err := json.Marshal(combinedRegisterData)
	if err != nil {
		log.Fatalf("Couldn't marshal JSON: %v", err)
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(combinedRegisterDataJSON))
	if err != nil {
		log.Fatalf("Error sending the registration request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading the response: %v", err)
	}

	client_secret := client_proof(body, ak)
	
	// send a subsequent request to the server with the secret to verify ourselves
	final_url := urlprefix + "/registration/finish"
	object_response := lib.ClientResponseData{Secret: client_secret}
	response_bytes, err := json.Marshal(object_response)

	if err != nil {
		log.Fatal("Unable to marshal response")
	}

	resp, err = client.Post(final_url, "application/json", bytes.NewBuffer(response_bytes))

	if err != nil {
		log.Fatalf("Error finishing the registration request: %v", err)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Couldn't decore the server's response")
	}

	fmt.Println(string(content))
}

func recieveMessage(urlprefix string, client *http.Client) {
	req, err := http.NewRequest("GET", urlprefix+"/download", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Not logged in")
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unkown error")
		return
	}

	var msgData lib.MessageData
	err = json.NewDecoder(resp.Body).Decode(&msgData)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Recieved messages:")
	fmt.Println(msgData.Content)
}

func sendMessage(urlprefix string, client *http.Client, messageContent string) error {
	msgData := lib.MessageData{
		Content: messageContent,
	}

	jsonData, err := json.Marshal(msgData)
	if err != nil {
		return fmt.Errorf("error marshalling message data: %v", err)
	}

	req, err := http.NewRequest("POST", urlprefix+"/message", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Message uploaded successfully.")
		return nil
	}

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Not logged in")
		return nil
	}

	fmt.Println(fmt.Errorf("Unknown error"))
	return nil
}

func PrintHelp() {
	fmt.Println()
	// general control
	fmt.Println("help - prints this message")
	fmt.Println("exit - halts this program")

	// registration
	fmt.Println("reqreg - requests registration")

	// message manipulation
	fmt.Println("pmess - prints current message to be sent")
	fmt.Println("chmess - changes current message to be sent")
	fmt.Println("smess - sends current message")

	// delete or TODO
	fmt.Println("dmess - download messages from server")

	// url manipulation
	fmt.Println("purl - prints current server url")
	fmt.Println("churl - changes currents server url")

	// ak manipulation
	fmt.Println("genkey - generates new ak")

	fmt.Println()
}

func UserPrompt(text string) string {
	var str string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, text)
		str, _ = r.ReadString('\n')
		if str != "" {
			break
		}
	}
	return strings.TrimSpace(str)
}

func main() {
	var url string = "https://localhost:8080" // base value
	var message string = "TestMessage 1234 %!"

	// create object for storing recieved cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal("Error creating cookie jar")
		return
	}

	// create a client which handles the reqeusts
	// this allows us to save the cookies and send them later
	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: os.Getenv("IGNORE_TLS_ERRORS") == "1"},
		},
	}

	// Generate EK and AK and initiate the registration process.
	ak, jsonData, err := client_generate_ek_ak()
	if err != nil {
		log.Fatalf("Error generating EK and AK: %v", err)
	}

	// dialog with user starts here
	fmt.Println("Hello, welcome to PV204 TPM project.\n")
	fmt.Println("Following commands are prepared:")
	PrintHelp()

	// main control body of programm
	for {
		command := UserPrompt("Enter a command: ")
		words := strings.Fields(command)

		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "":
			continue

		case "exit":
			return

		case "help":
			PrintHelp()

		case "purl":
			fmt.Println(url)

		case "churl":
			url = UserPrompt("Enter new url: ")

		case "pmess":
			fmt.Println(message)

		case "chmess":
			message = UserPrompt("Enter new message to be prepared: ")

		case "smess":
			err = sendMessage(url, client, message)
			if err != nil {
				log.Fatalf("Error sending message: %v", err)
			}

		case "dmess":
			recieveMessage(url, client)

		case "reqreg":
			if len(words) != 2 {
				PrintHelp()
			} else {
				authenticate(url, client, *jsonData, ak, words[1])
			}

		case "login":
			authenticate(url, client, *jsonData, ak, "")

		case "genkey":
			ak, jsonData, err = client_generate_ek_ak()
			if err != nil {
				log.Fatalf("Error generating EK and AK: %v", err)
			}

		default:
			fmt.Println("Unknown command\n")
		}
	}
}
