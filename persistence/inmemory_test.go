package persistence_test

import (
	"testing"

	"github.com/PaoloModica/signing-service-challenge-go/domain"
	test_utils "github.com/PaoloModica/signing-service-challenge-go/internal"
	"github.com/PaoloModica/signing-service-challenge-go/persistence"
)

func TestInMemorySignatureDeviceStore(t *testing.T) {
	t.Run("create InMemorySignatureDeviceStore", func(t *testing.T) {
		store, err := persistence.NewInMemorySignatureDeviceStore()
		test_utils.AssertErrorNotNil(t, "InMemorySignatureDeviceStore creation", err)

		devices, err := store.FindAll()
		test_utils.AssertErrorNotNil(t, "devices retrieval", err)
		expectedDevicesLen := 0
		if len(devices) != expectedDevicesLen {
			t.Errorf("expected devices length %d, found %d", expectedDevicesLen, len(devices))
		}

	})
	t.Run("InMemorySignatureDeviceStore capabilities", func(t *testing.T) {
		store, _ := persistence.NewInMemorySignatureDeviceStore()
		device, _ := domain.NewSignatureDevice("testDevice", []byte("privateKey"), domain.ECC)
		t.Run("create new signature device", func(t *testing.T) {
			devices, _ := store.FindAll()
			expectedDevicesLen := len(devices) + 1

			err := store.Create(device)
			test_utils.AssertErrorNotNil(t, "device creation and storaging", err)

			devices, _ = store.FindAll()
			test_utils.AssertSignatureDeviceStoreLen(t, expectedDevicesLen, len(devices))
		})
		t.Run("find signature device by ID, existing device", func(t *testing.T) {
			requestedDevice, err := store.FindById(device.Id)
			test_utils.AssertErrorNotNil(t, "device retrieval", err)
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
			test_utils.AssertErrorNotNil(t, "device update and storaging", err)

			updatedDevice, _ := store.FindById(device.Id)
			gotCounter := updatedDevice.GetSignatureCounter()

			if expectedCounter != gotCounter {
				t.Errorf("expected device signature counter to be %d, got %d", expectedCounter, gotCounter)
			}
		})
		t.Run("update signature device counter, unknown device", func(t *testing.T) {
			deviceNotInStore, _ := domain.NewSignatureDevice("newDevice", []byte("privateKey"), domain.ECC)

			err := store.Update(deviceNotInStore)
			if err == nil {
				t.Errorf("expected device not to be found")
			}
		})
	})
}
