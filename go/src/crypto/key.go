package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
)

var curve = elliptic.P256()

func GenPKey() (*ecdsa.PrivateKey, error) {
	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ECDSA key pair: %w", err)
	}

	return privateKey, nil
}

func PubKeyToBytes(pubKey ecdsa.PublicKey) []byte {
	return append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)
}

func PubKeyFromBytes(bytes []byte) *ecdsa.PublicKey {
	len := curve.Params().BitSize / 8
	return &ecdsa.PublicKey{
		Curve: curve,
		X:     new(big.Int).SetBytes(bytes[:len]),
		Y:     new(big.Int).SetBytes(bytes[len:]),
	}
}
