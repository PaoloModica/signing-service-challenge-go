package domain_test

import (
	"testing"

	"github.com/PaoloModica/signing-service-challenge-go/domain"
	test_utils "github.com/PaoloModica/signing-service-challenge-go/internal"
)

type signatureDeviceTestCase struct {
	description string
	label       string
	privateKey  []byte
	keytype     domain.KeyGenAlgorithm
}

func TestSignatureDevice(t *testing.T) {
	t.Run("create SignatureDevice instance", func(t *testing.T) {
		signatureDevicesTestCases := []signatureDeviceTestCase{
			{"RSA signature device", "myRSADevice", []byte("PrivateKey"), domain.RSA},
			{"ECC signature device", "myECCDevice", []byte("PublicKey"), domain.RSA},
		}
		for _, tc := range signatureDevicesTestCases {
			t.Run(tc.description, func(t *testing.T) {
				device, err := domain.NewSignatureDevice(tc.label, tc.privateKey, tc.keytype)

				test_utils.AssertErrorNotNil(t, "signature device creation", err)
				assertSignatureDeviceInitialStatus(t, device)
			})
		}
	})

	t.Run("set last signature and get SignatureDevice instance counter", func(t *testing.T) {
		device, err := domain.NewSignatureDevice("device", []byte("privateKey"), domain.RSA)

		lastSignature := []byte("lastSignature")
		test_utils.AssertErrorNotNil(t, "signature device creation", err)

		device.SetLastSignature(lastSignature)
		expectedVal := 1
		gotVal := device.GetSignatureCounter()

		if expectedVal != gotVal {
			t.Errorf("expected signature device counter to be 1, found")
		}

		gotSignature, _ := device.GetLastSignature()
		if string(gotSignature) != string(lastSignature) {
			t.Errorf("expected signature device last signature to be %s, found %s", string(lastSignature), string(gotSignature))
		}
	})
}

func TestSignatureDeviceRepository(t *testing.T) {
	store := test_utils.StubSignatureDeviceStore{
		Store: map[string]*domain.SignatureDevice{},
	}
	t.Run("create new signature device repository", func(t *testing.T) {
		repository, err := domain.NewSignatureDeviceRepository(&store)

		if repository == nil || err != nil {
			t.Errorf("expected SignatureDeviceRepository to have been created")
		}
	})
	t.Run("SignatureDeviceRepository capabilities", func(t *testing.T) {
		repository, _ := domain.NewSignatureDeviceRepository(&store)
		device, _ := domain.NewSignatureDevice("testDevice", []byte("privateKey"), domain.RSA)

		t.Run("create new signature device", func(t *testing.T) {
			devices, _ := repository.FindAll()
			expectedDevicesLen := len(devices) + 1

			id, err := repository.Create(device)
			test_utils.AssertSignatureDeviceId(t, id, err)

			devices, _ = repository.FindAll()
			test_utils.AssertSignatureDeviceStoreLen(t, expectedDevicesLen, len(devices))
		})
		t.Run("find device by ID, existing device", func(t *testing.T) {
			requestedDevice, err := repository.FindById(device.Id)
			test_utils.AssertErrorNotNil(t, "device retrieval", err)
			if requestedDevice != device {
				t.Errorf("expected %s device, found %s", device.Label, requestedDevice.Label)
			}
		})
		t.Run("update signature device last signature, existing device", func(t *testing.T) {
			lastSignature := []byte("lastSignature")
			expectedCounter := device.GetSignatureCounter() + 1
			device.SetLastSignature(lastSignature)

			err := repository.Update(device)
			test_utils.AssertErrorNotNil(t, "device update and storaging", err)

			updatedDevice, _ := store.FindById(device.Id)
			gotCounter := updatedDevice.GetSignatureCounter()

			if expectedCounter != gotCounter {
				t.Errorf("expected device signature counter to be %d, got %d", expectedCounter, gotCounter)
			}
		})
	})
}

func TestSignatureDeviceService(t *testing.T) {
	device, _ := domain.NewSignatureDevice("testDevice1", []byte("privateKey"), "RSA")
	store := test_utils.StubSignatureDeviceStore{
		Store: map[string]*domain.SignatureDevice{device.Id: device},
	}
	repository, _ := domain.NewSignatureDeviceRepository(&store)
	t.Run("create new signature device service", func(t *testing.T) {
		service, err := domain.NewSignatureDeviceService(repository)

		if service == nil || err != nil {
			t.Errorf("expected SignatureDeviceRepository to have been created")
		}
	})
	t.Run("signature device service capabilities", func(t *testing.T) {
		service, _ := domain.NewSignatureDeviceService(repository)
		t.Run("find all signature devices", func(t *testing.T) {
			expectedDeviceLen := 1
			devices, err := service.FindAll()
			test_utils.AssertErrorNotNil(t, "device retrieval", err)
			test_utils.AssertSignatureDeviceStoreLen(t, expectedDeviceLen, len(devices))
		})
		t.Run("find signature device by ID", func(t *testing.T) {
			d, err := service.FindById(device.Id)
			test_utils.AssertErrorNotNil(t, "device retrieval", err)
			if d != device {
				t.Errorf("expected to find device %s, got device %s", device.Label, d.Label)
			}
		})
		t.Run("create new signature device", func(t *testing.T) {
			devices, _ := service.FindAll()
			expectedDeviceLen := len(devices) + 1

			id, err := service.Create("testDevice", "RSA")
			test_utils.AssertSignatureDeviceId(t, id, err)

			devices, _ = service.FindAll()
			test_utils.AssertSignatureDeviceStoreLen(t, expectedDeviceLen, len(devices))
		})
		t.Run("update signature device counter, existing device", func(t *testing.T) {
			lastSignature := []byte("lastSignature")
			expectedCounter := device.GetSignatureCounter() + 1

			err := service.Update(device.Id, lastSignature)
			test_utils.AssertErrorNotNil(t, "device update and storaging", err)

			updatedDevice, _ := service.FindById(device.Id)
			gotCounter := updatedDevice.GetSignatureCounter()

			if expectedCounter != gotCounter {
				t.Errorf("expected device signature counter to be %d, got %d", expectedCounter, gotCounter)
			}
		})
		t.Run("update device with unknown ID", func(t *testing.T) {
			lastSignature := []byte("lastSignature")
			err := service.Update("unknownId", lastSignature)
			if err == nil {
				t.Errorf("expected not found error")
			}
		})
	})
}

func assertSignatureDeviceInitialStatus(t *testing.T, d *domain.SignatureDevice) {
	t.Helper()

	if d.Id == "" {
		t.Errorf("expected signature device ID to be set, found empty")
	}

	deviceCounter := d.GetSignatureCounter()
	if deviceCounter > 0 {
		t.Errorf("expected signature device counter to be initialised at 0, found %d", deviceCounter)
	}
}
