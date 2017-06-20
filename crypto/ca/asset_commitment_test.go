package ca

import (
	"encoding/hex"
	"testing"

	"chain/crypto/ed25519/ecmath"
)

func TestAssetCommitment(t *testing.T) {
	var assetID AssetID
	hex.Decode(assetID[:], []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))

	ac1, c := CreateAssetCommitment(assetID, nil)
	if !ac1.C().ConstTimeEqual(&ecmath.ZeroPoint) {
		t.Error("expected zero point")
	}
	if c != nil {
		t.Error("expected nil blinding factor")
	}

	var aekBuf [32]byte
	hex.Decode(aekBuf[:], []byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"))

	ac2, c := CreateAssetCommitment(assetID, aekBuf[:])
	if ac2.C().ConstTimeEqual(&ecmath.ZeroPoint) {
		t.Error("expected nonzero point")
	}
	if ac1.H().ConstTimeEqual(ac2.H()) {
		t.Error("expected different H value")
	}
	if c == nil {
		t.Error("expected non-nil blinding factor")
	}
}
