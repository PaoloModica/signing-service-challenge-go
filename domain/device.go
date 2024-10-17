package domain

import (
	"github.com/google/uuid"
)

type SignatureDevice struct {
	Id               string
	Label            string
	PrivateKey       []byte
	PublicKey        []byte
	signatureCounter int
}

func NewSignatureDevice(label string, privateKey []byte, publicKey []byte) (*SignatureDevice, error) {
	return &SignatureDevice{Id: uuid.NewString(), Label: label, PrivateKey: privateKey, PublicKey: publicKey}, nil
}

func (s *SignatureDevice) GetSignatureCounter() int {
	return s.signatureCounter
}

func (s *SignatureDevice) IncrementSignatureCounter() {
	s.signatureCounter += 1
}
