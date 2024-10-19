package test_utils

import (
	"testing"
)

func AssertErrorNotNil(t *testing.T, message string, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("An error occurred during %s, error: %s", message, err.Error())
	}
}

func AssertSignatureDeviceStoreLen(t *testing.T, exp int, got int) {
	t.Helper()

	if exp != got {
		t.Errorf("expected %d devices in store, got %d", exp, got)
	}
}
