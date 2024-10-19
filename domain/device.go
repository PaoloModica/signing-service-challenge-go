package domain

import (
	"fmt"

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

type SignatureDeviceStore interface {
	FindById(id string) (*SignatureDevice, error)
	FindAll() ([]*SignatureDevice, error)
	Create(*SignatureDevice) error
	Update(*SignatureDevice) error
}

type DeviceNotFoundError struct {
	SignatureDeviceId string
}

func (e *DeviceNotFoundError) Error() string {
	return fmt.Sprintf("device with ID %s not found", e.SignatureDeviceId)
}
