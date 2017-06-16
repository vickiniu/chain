package ca

import (
	"chain/crypto/ed25519/ecmath"
	"testing"
)

func TestEncryptValue(t *testing.T) {
	var assetID AssetID
	value := uint64(17)
	vek := []byte("value encryption key")

	ac, _ := CreateAssetCommitment(assetID, nil)
	vc, f := CreateValueCommitment(value, ac, vek)

	evef := make([]byte, EncryptedValueSize)
	EncryptValue(vc, value, f, vek, evef)

	value2, f2, ok := DecryptValue(evef, vc, ac, vek)

	if !ok {
		t.Error("decryption failed")
	} else {
		if value != value2 {
			t.Errorf("got value %d, want %d", value2, value)
		}
		if !f.Equal(f2) {
			t.Errorf("got blinding factor %x, want %x", f2[:], f[:])
		}
	}

	for i := 0; i < len(evef); i++ {
		for j := uint(0); j < 8; j++ {
			evef[i] ^= 1 << j

			value2, f2, ok := DecryptValue(evef, vc, ac, vek)

			if ok {
				t.Error("unexpected decryption success with bad encrypted value amount")
			}

			if value2 != 0 {
				t.Error("unexpected value from failed decryption")
			}

			if !f2.Equal(&ecmath.Zero) {
				t.Error("unexpected value from failed decryption")
			}

			evef[i] ^= 1 << j
		}
	}
}