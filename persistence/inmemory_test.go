package persistence_test

import (
	"testing"

	"github.com/PaoloModica/signing-service-challenge-go/domain"
	"github.com/PaoloModica/signing-service-challenge-go/persistence"
)

func TestInMemorySignatureDeviceStore(t *testing.T) {
	t.Run("create InMemorySignatureDeviceStore", func(t *testing.T) {
		store, err := persistence.NewInMemorySignatureDeviceStore()
		assertInMemorySignatureDeviceStoreError(t, err)

		devices, err := store.FindAll()
		assertInMemorySignatureDeviceStoreError(t, err)
		expectedDevicesLen := 0
		if len(devices) != expectedDevicesLen {
			t.Errorf("expected devices length %d, found %d", expectedDevicesLen, len(devices))
		}

	})
	t.Run("InMemorySignatureDeviceStore capabilities", func(t *testing.T) {
		store, _ := persistence.NewInMemorySignatureDeviceStore()
		device, _ := domain.NewSignatureDevice("testDevice", []byte("privateKey"), []byte("publicKey"))
		t.Run("create new signature device", func(t *testing.T) {
			devices, _ := store.FindAll()
			expectedDevicesLen := len(devices) + 1

			err := store.Create(device)
			assertInMemorySignatureDeviceStoreError(t, err)

			devices, _ = store.FindAll()
			assertSignatureDeviceStoreLen(t, expectedDevicesLen, len(devices))
		})
		t.Run("find signature device by ID, existing device", func(t *testing.T) {
			requestedDevice, err := store.FindById(device.Id)
			assertInMemorySignatureDeviceStoreError(t, err)
			if requestedDevice != device {
				t.Errorf("expected %s device, found %s", device.Label, requestedDevice.Label)
			}
		})
		t.Run("find signature device by ID, unknown device", func(t *testing.T) {
			_, err := store.FindById("unknownDevice")
			if err == nil {
				t.Errorf("expected device not to be found")
			}
		})
		t.Run("update signature device counter, existing device", func(t *testing.T) {
			expectedCounter := device.GetSignatureCounter() + 1
			device.IncrementSignatureCounter()

			err := store.Update(device)
			assertInMemorySignatureDeviceStoreError(t, err)

			updatedDevice, _ := store.FindById(device.Id)
			gotCounter := updatedDevice.GetSignatureCounter()

			if expectedCounter != gotCounter {
				t.Errorf("expected device signature counter to be %d, got %d", expectedCounter, gotCounter)
			}
		})
		t.Run("update signature device counter, unknown device", func(t *testing.T) {
			deviceNotInStore, _ := domain.NewSignatureDevice("newDevice", []byte("privateKey"), []byte("publicKey"))

			err := store.Update(deviceNotInStore)
			if err == nil {
				t.Errorf("expected device not to be found")
			}
		})
	})
}

func assertInMemorySignatureDeviceStoreError(t *testing.T, e error) {
	t.Helper()

	if e != nil {
		t.Errorf("an error occurred: %s", e)
	}
}

func assertSignatureDeviceStoreLen(t *testing.T, exp int, got int) {
	t.Helper()

	if exp != got {
		t.Errorf("expected %d devices in store, got %d", exp, got)
	}
}
