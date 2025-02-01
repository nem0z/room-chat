package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

func GenWitness(pKey *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)

	r, s, err := ecdsa.Sign(rand.Reader, pKey, hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign the hash: %v", err)
	}

	return append(r.Bytes(), s.Bytes()...), nil
}

func VerifyWitness(pubKeyBytes []byte, witness []byte, data []byte) bool {
	hash := sha256.Sum256(data)
	pubKey := PubKeyFromBytes(pubKeyBytes)

	rLen := pubKey.Curve.Params().BitSize / 8
	r := new(big.Int).SetBytes(witness[:rLen])
	s := new(big.Int).SetBytes(witness[rLen:])

	return ecdsa.Verify(pubKey, hash[:], r, s)
}
