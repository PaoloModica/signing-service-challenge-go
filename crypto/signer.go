package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// TODO: implement RSA and ECDSA signing ...
type RSASigner struct {
	devicePrivateKey []byte
	lastSignature    string
	signatureCount   int
	marshaler        *RSAMarshaler
}

func NewRSASigner(devicePrivateKey []byte, lastSignature string, signatureCount int) (*RSASigner, error) {
	return &RSASigner{devicePrivateKey: devicePrivateKey, lastSignature: lastSignature, signatureCount: signatureCount, marshaler: &RSAMarshaler{}}, nil
}

func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	keyPair, err := s.marshaler.Unmarshal(s.devicePrivateKey)
	if err != nil {
		log.Fatalf("an error occurred while unmarshalling private key: %s", err.Error())
		return nil, err
	}
	encodedLastSignature := base64.StdEncoding.EncodeToString([]byte(s.lastSignature))
	signatureInput := fmt.Sprintf("%d_%s_%s", s.signatureCount, string(dataToBeSigned), encodedLastSignature)

	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(signatureInput))
	if err != nil {
		log.Fatalf("an error occurred while hashing data to be signed: %s", err.Error())
	}
	msgHashSum := msgHash.Sum(nil)

	return rsa.SignPSS(rand.Reader, keyPair.Private, crypto.SHA256, msgHashSum, nil)
}

type ECDSASigner struct {
	devicePrivateKey []byte
	lastSignature    string
	signatureCount   int
	marshaler        *ECCMarshaler
}

func NewECDSASigner(devicePrivateKey []byte, lastSignature string, signatureCount int) (*ECDSASigner, error) {
	return &ECDSASigner{devicePrivateKey: devicePrivateKey, lastSignature: lastSignature, signatureCount: signatureCount, marshaler: &ECCMarshaler{}}, nil
}

func (s *ECDSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	keyPair, err := s.marshaler.Decode(s.devicePrivateKey)
	if err != nil {
		log.Fatalf("an error occurred while unmarshalling private key: %s", err.Error())
		return nil, err
	}
	encodedLastSignature := base64.StdEncoding.EncodeToString([]byte(s.lastSignature))
	signatureInput := fmt.Sprintf("%d_%s_%s", s.signatureCount, string(dataToBeSigned), encodedLastSignature)

	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(signatureInput))
	if err != nil {
		log.Fatalf("an error occurred while hashing data to be signed: %s", err.Error())
	}
	msgHashSum := msgHash.Sum(nil)

	return ecdsa.SignASN1(rand.Reader, keyPair.Private, msgHashSum)
}
