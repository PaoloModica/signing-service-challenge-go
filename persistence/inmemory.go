package persistence

import (
	"fmt"
	"log"

	"github.com/PaoloModica/signing-service-challenge-go/domain"
)

type InMemorySignatureDeviceStore struct {
	store map[string]*domain.SignatureDevice
}

func NewInMemorySignatureDeviceStore() (*InMemorySignatureDeviceStore, error) {
	return &InMemorySignatureDeviceStore{map[string]*domain.SignatureDevice{}}, nil
}

func (s *InMemorySignatureDeviceStore) FindById(id string) (*domain.SignatureDevice, error) {
	device, found := s.store[id]

	if !found {
		return nil, domain.DeviceNotFoundError(fmt.Sprintf("device with ID %s not found", id))
	}
	return device, nil
}

func (s *InMemorySignatureDeviceStore) FindAll() ([]*domain.SignatureDevice, error) {
	devices := []*domain.SignatureDevice{}
	for _, v := range s.store {
		devices = append(devices, v)
	}
	return devices, nil
}

func (s *InMemorySignatureDeviceStore) Create(d *domain.SignatureDevice) (string, error) {
	s.store[d.Id] = d
	log.Printf("device %s stored successfully", d.Label)
	return d.Id, nil
}

func (s *InMemorySignatureDeviceStore) Update(d *domain.SignatureDevice) error {
	_, found := s.store[d.Id]

	if !found {
		return domain.DeviceNotFoundError(fmt.Sprintf("device with ID %s not found", d.Id))
	}

	s.store[d.Id] = d
	log.Printf("device %s updated successfully", d.Label)
	return nil
}
