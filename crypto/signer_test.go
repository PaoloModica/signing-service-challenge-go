package crypto_test

import (
	"testing"

	"github.com/PaoloModica/signing-service-challenge-go/crypto"
	"github.com/google/uuid"
)

func TestSigner(t *testing.T) {
	RSAMarshaler := crypto.NewRSAMarshaler()
	RSAKeyGen := &crypto.RSAGenerator{}
	RSAPrivateKey, _ := RSAKeyGen.Generate()
	_, marshalledRSAPrivateKey, _ := RSAMarshaler.Marshal(*RSAPrivateKey)

	ECCMarshaler := crypto.NewECCMarshaler()
	ECCKeyGen := &crypto.ECCGenerator{}
	ECCPrivateKey, _ := ECCKeyGen.Generate()
	_, marshalledECCPrivateKey, _ := ECCMarshaler.Encode(*ECCPrivateKey)

	signerParams := struct {
		devicePrivateKey []byte
		lastSignature    string
		signatureCount   int
	}{
		devicePrivateKey: []byte("privateKey"),
		lastSignature:    uuid.NewString(),
		signatureCount:   0,
	}
	t.Run("create new RSA signer", func(t *testing.T) {
		signer, err := crypto.NewRSASigner(signerParams.devicePrivateKey, signerParams.lastSignature, signerParams.signatureCount)
		if signer == nil || err != nil {
			t.Errorf("expected RSA signer to be created, got error: %s", err.Error())
		}
	})
	t.Run("create new ECDSA signer", func(t *testing.T) {
		signer, err := crypto.NewECDSASigner(signerParams.devicePrivateKey, signerParams.lastSignature, signerParams.signatureCount)
		if signer == nil || err != nil {
			t.Errorf("expected RSA signer to be created, got error: %s", err.Error())
		}
	})
	t.Run("sign data", func(t *testing.T) {
		rsaSigner, _ := crypto.NewRSASigner(marshalledRSAPrivateKey, signerParams.lastSignature, signerParams.signatureCount)
		ecdsaSigner, _ := crypto.NewECDSASigner(marshalledECCPrivateKey, signerParams.lastSignature, signerParams.signatureCount)
		dataToBeSigned := []byte("test data")

		t.Run("sign with RSA", func(t *testing.T) {
			signature, err := rsaSigner.Sign(dataToBeSigned)

			if signature == nil || err != nil {
				t.Errorf("expected RSA signature to be created, got none")
			}
		})
		t.Run("sign with ECDSA", func(t *testing.T) {
			signature, err := ecdsaSigner.Sign(dataToBeSigned)

			if signature == nil || err != nil {
				t.Errorf("expected ECDSA signature to be created, got none")
			}
		})
	})
}
