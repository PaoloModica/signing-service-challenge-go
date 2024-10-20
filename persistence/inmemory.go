package persistence

import (
	"fmt"
	"log"
	"sync"

	"github.com/PaoloModica/signing-service-challenge-go/domain"
)

type InMemorySignatureDeviceStore struct {
	store map[string]*domain.SignatureDevice
	lock  sync.RWMutex
}

func NewInMemorySignatureDeviceStore() (*InMemorySignatureDeviceStore, error) {
	return &InMemorySignatureDeviceStore{map[string]*domain.SignatureDevice{}, sync.RWMutex{}}, nil
}

func (s *InMemorySignatureDeviceStore) FindById(id string) (*domain.SignatureDevice, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

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
	s.lock.Lock()
	defer s.lock.Unlock()

	s.store[d.Id] = d
	log.Printf("device %s stored successfully", d.Label)
	return d.Id, nil
}

func (s *InMemorySignatureDeviceStore) Update(d *domain.SignatureDevice) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, found := s.store[d.Id]

	if !found {
		return domain.DeviceNotFoundError(fmt.Sprintf("device with ID %s not found", d.Id))
	}

	s.store[d.Id] = d
	log.Printf("device %s updated successfully", d.Label)
	return nil
}
