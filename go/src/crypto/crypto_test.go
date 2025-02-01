package crypto_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/nem0z/room-chat/src/crypto"
)

func TestGenWitnessAndVerifyWitness(t *testing.T) {
	// Generate private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Test data
	data := []byte("test data")

	// Generate the witness
	witness, err := crypto.GenWitness(privateKey, data)
	if err != nil {
		t.Fatalf("GenWitness failed: %v", err)
	}

	// Verify the witness
	pubKeyBytes := crypto.PubKeyToBytes(privateKey.PublicKey)
	if !crypto.VerifyWitness(pubKeyBytes, witness, data) {
		t.Fatalf("VerifyWitness failed")
	}

	// Test with tampered data
	tamperedData := []byte("tampered data")
	if crypto.VerifyWitness(pubKeyBytes, witness, tamperedData) {
		t.Fatalf("VerifyWitness should have failed with tampered data")
	}

	// Test with tampered witness
	tamperedWitness := make([]byte, len(witness))
	copy(tamperedWitness, witness)
	tamperedWitness[0] ^= 0xFF // Flip some bits in the witness
	if crypto.VerifyWitness(pubKeyBytes, tamperedWitness, data) {
		t.Fatalf("VerifyWitness should have failed with tampered witness")
	}

	otherWitness, err := crypto.GenWitness(privateKey, tamperedData)
	if err != nil {
		t.Fatalf("GenWitness failed: %v", err)
	}

	if crypto.VerifyWitness(pubKeyBytes, otherWitness, data) {
		t.Fatalf("VerifyWitness should have failed with another witness from the same public")
	}
}

func TestGenWitness_ErrorHandling(t *testing.T) {
	// Generate a new ECDSA private key
	pKey, err := crypto.GenPKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Test with nil private key
	_, err = crypto.GenWitness(nil, []byte("test data"))
	if err == nil {
		t.Fatalf("GenWitness should have failed with nil private key")
	}

	// Test with nil data
	_, err = crypto.GenWitness(pKey, nil)
	if err != nil {
		t.Fatalf("GenWitness should handle nil data without error")
	}
}

func TestVerifyWitness_ErrorHandling(t *testing.T) {
	// Generate a new ECDSA private key
	pKey, err := crypto.GenPKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Test data
	data := []byte("test data")

	// Generate the witness
	witness, err := crypto.GenWitness(pKey, data)
	if err != nil {
		t.Fatalf("GenWitness failed: %v", err)
	}

	// Test with nil public key bytes
	if crypto.VerifyWitness(nil, witness, data) {
		t.Fatalf("VerifyWitness should have failed with nil public key bytes")
	}

	// Test with invalid public key bytes
	invalidPubKeyBytes := []byte("invalid public key bytes")
	if crypto.VerifyWitness(invalidPubKeyBytes, witness, data) {
		t.Fatalf("VerifyWitness should have failed with invalid public key bytes")
	}

	// Test with nil witness
	pubKeyBytes := crypto.PubKeyToBytes(pKey.PublicKey)
	if crypto.VerifyWitness(pubKeyBytes, nil, data) {
		t.Fatalf("VerifyWitness should have failed with nil witness")
	}

	witness, err = crypto.GenWitness(pKey, nil)
	if err != nil {
		t.Fatalf("GenWitness failed: %v", err)
	}

	// Test with nil data
	if !crypto.VerifyWitness(pubKeyBytes, witness, nil) {
		t.Fatalf("VerifyWitness should handle nil data without error")
	}
}
