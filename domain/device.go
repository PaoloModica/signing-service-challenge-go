package domain

import (
	"fmt"
	"log"
	"sync"

	"github.com/PaoloModica/signing-service-challenge-go/crypto"
	"github.com/google/uuid"
)

type KeyGenAlgorithm string

const (
	ECC KeyGenAlgorithm = "ECC"
	RSA KeyGenAlgorithm = "RSA"
)

type KeyTypeNotValidError string

func (e KeyTypeNotValidError) Error() string {
	return string(e)
}

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
	Create(*SignatureDevice) (string, error)
	Update(*SignatureDevice) error
}

type DeviceNotFoundError string

func (e DeviceNotFoundError) Error() string {
	return string(e)
}

type SignatureDeviceRepository interface {
	FindById(id string) (*SignatureDevice, error)
	FindAll() ([]*SignatureDevice, error)
	Create(*SignatureDevice) (string, error)
	Update(*SignatureDevice) error
}

type signatureDeviceRepository struct {
	lock  sync.RWMutex
	store SignatureDeviceStore
}

func NewSignatureDeviceRepository(s SignatureDeviceStore) (*signatureDeviceRepository, error) {
	return &signatureDeviceRepository{lock: sync.RWMutex{}, store: s}, nil
}

func (r *signatureDeviceRepository) FindById(id string) (*SignatureDevice, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.store.FindById(id)
}

func (r *signatureDeviceRepository) FindAll() ([]*SignatureDevice, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.store.FindAll()
}

func (r *signatureDeviceRepository) Create(d *SignatureDevice) (string, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.store.Create(d)
}

func (r *signatureDeviceRepository) Update(d *SignatureDevice) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.store.Update(d)
}

type SignatureDeviceService interface {
	FindById(id string) (*SignatureDevice, error)
	FindAll() ([]*SignatureDevice, error)
	Create(label string, keyType KeyGenAlgorithm) (string, error)
	Update(id string) error
}

type signatureDeviceService struct {
	repository      SignatureDeviceRepository
	rsaKeyGenerator *crypto.RSAGenerator
	rsaKeyMarshaler *crypto.RSAMarshaler
	eccKeyGenerator *crypto.ECCGenerator
	eccKeyMarshaler *crypto.ECCMarshaler
}

func NewSignatureDeviceService(repository SignatureDeviceRepository) (*signatureDeviceService, error) {
	return &signatureDeviceService{repository: repository, rsaKeyGenerator: &crypto.RSAGenerator{}, rsaKeyMarshaler: &crypto.RSAMarshaler{}, eccKeyGenerator: &crypto.ECCGenerator{}, eccKeyMarshaler: &crypto.ECCMarshaler{}}, nil
}

func (s *signatureDeviceService) FindAll() ([]*SignatureDevice, error) {
	return s.repository.FindAll()
}

func (s *signatureDeviceService) FindById(id string) (*SignatureDevice, error) {
	return s.repository.FindById(id)
}

func (s *signatureDeviceService) createAndDecodePrivateKey(keyType KeyGenAlgorithm) ([]byte, error) {
	switch keyType {
	case RSA:
		keyPair, err := s.rsaKeyGenerator.Generate()
		if err != nil {
			log.Fatalf("an error occurred while generating RSA key pair: %s", err.Error())
			return nil, err
		}
		_, privateKey, err := s.rsaKeyMarshaler.Marshal(*keyPair)
		if err != nil {
			log.Fatalf("an error occurred while marshaling RSA key pair: %s", err.Error())
			return nil, err
		} else {
			return privateKey, nil
		}
	case ECC:
		keyPair, err := s.eccKeyGenerator.Generate()
		if err != nil {
			log.Fatalf("an error occurred while generating ECC key pair: %s", err.Error())
			return nil, err
		}
		_, privateKey, err := s.eccKeyMarshaler.Encode(*keyPair)
		if err != nil {
			log.Fatalf("an error occurred while marshaling ECC key pair: %s", err.Error())
			return nil, err
		} else {
			return privateKey, nil
		}
	default:
		return nil, KeyTypeNotValidError("key generation algorithm not valid or unknown")
	}
}

func (s *signatureDeviceService) Create(label string, keyType KeyGenAlgorithm) (string, error) {
	privateKey, err := s.createAndDecodePrivateKey(keyType)
	if err != nil {
		log.Fatalf("an error occurred while creating signature device: %s", err.Error())
		return "", err
	}
	device, err := NewSignatureDevice(label, privateKey, keyType)
	if err != nil {
		log.Fatalf("an error occurred while creating signature device: %s", err.Error())
		return "", err
	}
	id, err := s.repository.Create(device)
	if err != nil {
		log.Fatalf("an error occurred while creating signature device: %s", err.Error())
		return "", err
	}
	return id, nil
}

func (s *signatureDeviceService) Update(id string) error {
	device, err := s.repository.FindById(id)
	if device == nil || err != nil {
		return DeviceNotFoundError(fmt.Sprintf("device with ID %s not found", id))
	}
	device.IncrementSignatureCounter()
	return s.repository.Update(device)
}
