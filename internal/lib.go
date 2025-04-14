package lib

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rsa"
	"log"

	"github.com/google/go-attestation/attest"
)

const CookieName = "session"

// PublicKeyType defines the types of EK public keys.
type PublicKeyType int

const (
	RSA   PublicKeyType = iota // RSA key type
	ECDH                       // ECDH key type
	ECDSA                      // ECDSA key type
)

// RegisterData is the data sent by the client to the server the first time they register.
type RegisterData struct {
	Version      attest.TPMVersion            // The version of the TPM.
	PublicKey    string                       // The public key of the EK.
	AttestParams attest.AttestationParameters // The attestation parameters.
}

type MessageData struct {
	Content string `json:"content"`
}

// PreRegisterData is the data sent by the client to the server before registration
// to tell the server the type of the public key.
type PreRegisterData struct {
	PublicKeyType PublicKeyType
}

// RSARegisterData specializes RegisterData with an RSA public key.
type RSARegisterData struct {
	RegisterData
	PublicKey rsa.PublicKey
}

// ECDHRegisterData specializes RegisterData with an ECDH public key.
type ECDHRegisterData struct {
	RegisterData
	PublicKey ecdh.PublicKey
}

// ECDSARegisterData specializes RegisterData with an ECDSA public key.
type ECDSARegisterData struct {
	RegisterData
	PublicKey ecdsa.PublicKey
}

type CombinedRegisterData struct {
	PreData  PreRegisterData
	Data     RegisterData
	Username string
}

// ClientResponseData represents what the client responds to the server to prove its identity.
type ClientResponseData struct {
	Secret []byte
}

// handle_error is a utility function to handle errors.
func Handle_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
