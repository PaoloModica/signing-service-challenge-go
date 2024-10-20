package test_utils

import (
	"testing"
)

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
