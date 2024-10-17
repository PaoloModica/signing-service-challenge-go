package domain_test

import (
	"testing"

	"github.com/PaoloModica/signing-service-challenge-go/domain"
)

type signatureDeviceTestCase struct {
	description string
	label       string
	privateKey  []byte
	publicKey   []byte
}

func TestSignatureDevice(t *testing.T) {
	t.Run("create SignatureDevice instance", func(t *testing.T) {
		signatureDevicesTestCases := []signatureDeviceTestCase{
			{"RSA signature device", "myRSADevice", []byte("PublicKey"), []byte("PrivateKey")},
			{"ECC signature device", "myECCDevice", []byte("PublicKey"), []byte("PrivateKey")},
		}
		for _, tc := range signatureDevicesTestCases {
			t.Run(tc.description, func(t *testing.T) {
				device, err := domain.NewSignatureDevice(tc.label, tc.privateKey, tc.publicKey)

				assertSignatureDeviceCreationError(t, err)
				assertSignatureDeviceInitialStatus(t, device)
			})
		}
	})

	t.Run("increment and get SignatureDevice instance counter", func(t *testing.T) {
		device, err := domain.NewSignatureDevice("device", []byte("privateKey"), []byte("publicKey"))

		assertSignatureDeviceCreationError(t, err)

		device.IncrementSignatureCounter()
		expectedVal := 1
		gotVal := device.GetSignatureCounter()

		if expectedVal != gotVal {
			t.Errorf("expected signature device counter to be 1, found")
		}
	})
}

func assertSignatureDeviceCreationError(t *testing.T, e error) {
	t.Helper()

	if e != nil {
		t.Errorf("an error occurred while creating signature device: %s", e)
	}
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
