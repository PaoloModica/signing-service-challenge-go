package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type KeyGenAlgorithm string

const (
	ECC KeyGenAlgorithm = "ECC"
	RSA KeyGenAlgorithm = "RSA"
)

type SignatureDevice struct {
	Id               string
	Label            string
	PrivateKey       []byte
	KeyType          KeyGenAlgorithm
	signatureCounter int
}

func NewSignatureDevice(label string, privateKey []byte, keytype KeyGenAlgorithm) (*SignatureDevice, error) {
	return &SignatureDevice{Id: uuid.NewString(), Label: label, PrivateKey: privateKey, KeyType: keytype}, nil
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

type SignatureDeviceRepository interface {
	FindById(id string) (*SignatureDevice, error)
	FindAll() ([]*SignatureDevice, error)
	Create(*SignatureDevice) error
	Update(*SignatureDevice) error
}

type signatureDeviceRepository struct {
	store SignatureDeviceStore
}

func NewSignatureDeviceRepository(s SignatureDeviceStore) (*signatureDeviceRepository, error) {
	return &signatureDeviceRepository{store: s}, nil
}

func (r *signatureDeviceRepository) FindById(id string) (*SignatureDevice, error) {
	return r.store.FindById(id)
}

func (r *signatureDeviceRepository) FindAll() ([]*SignatureDevice, error) {
	return r.store.FindAll()
}

func (r *signatureDeviceRepository) Create(d *SignatureDevice) error {
	return r.store.Create(d)
}

func (r *signatureDeviceRepository) Update(d *SignatureDevice) error {
	return r.store.Update(d)
}
