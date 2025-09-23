package crypt

import (
	"crypto/rand"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
)

func LoadRSAKey(fileName string) (crypto.PrivKey, error) {
	privBytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	priv, err := crypto.UnmarshalPrivateKey(privBytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func SaveRSAKey(fileName string, priv crypto.PrivKey) error {
	privBytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, privBytes, 0600)
	if err != nil {
		return err
	}

	return nil
}

func GenRSAKey() (crypto.PrivKey, error) {
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)

	return priv, err
}
