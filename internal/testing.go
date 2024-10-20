package test_utils

import (
	"fmt"
	"testing"

	"github.com/PaoloModica/signing-service-challenge-go/domain"
)

type StubSignatureDeviceStore struct {
	Store map[string]*domain.SignatureDevice
}

func (s *StubSignatureDeviceStore) FindById(id string) (*domain.SignatureDevice, error) {
	d, found := s.Store[id]
	if !found {
		return nil, domain.DeviceNotFoundError(fmt.Sprintf("device with ID %s not found", id))
	}
	return d, nil
}

func (s *StubSignatureDeviceStore) FindAll() ([]*domain.SignatureDevice, error) {
	devices := []*domain.SignatureDevice{}
	for _, v := range s.Store {
		devices = append(devices, v)
	}
	return devices, nil
}

func (s *StubSignatureDeviceStore) Create(d *domain.SignatureDevice) (string, error) {
	s.Store[d.Id] = d
	return d.Id, nil
}

func (s *StubSignatureDeviceStore) Update(d *domain.SignatureDevice) error {
	d, err := s.FindById(d.Id)
	if err != nil {
		return nil
	}
	s.Store[d.Id] = d
	return nil
}

func AssertErrorNotNil(t *testing.T, message string, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("an error occurred during %s, error: %s", message, err.Error())
	}
}

func AssertSignatureDeviceId(t *testing.T, id string, err error) {
	t.Helper()

	if id == "" || err != nil {
		t.Error("expected signature device ID not null")
	}
}

func AssertSignatureDeviceStoreLen(t *testing.T, exp int, got int) {
	t.Helper()

	if exp != got {
		t.Errorf("expected %d devices in store, got %d", exp, got)
	}
}
