package crypt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"
	blst "github.com/supranational/blst/bindings/go"
)

// LibP2P RSA Private key operations (we don't care about public keys)
func GenRSAKey() (crypto.PrivKey, error) {
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)

	return priv, err
}

func Base64RSAKey(priv crypto.PrivKey) (string, error) {
	privBytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return "", err
	}

	s := base64.StdEncoding.EncodeToString(privBytes)

	return s, nil
}

func Unbase64RSAKey(b64rsa string) (crypto.PrivKey, error) {
	privBytes, err := base64.StdEncoding.DecodeString(b64rsa)
	if err != nil {
		return nil, err
	}

	priv, err := crypto.UnmarshalPrivateKey(privBytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// BLS Secret Key Operations
func GenBlsSk() (*blst.SecretKey, error) {
	var ikm [32]byte
	_, err := rand.Read(ikm[:])
	if err != nil {
		return nil, err
	}

	sk := blst.KeyGen(ikm[:])
	if sk == nil {
		return nil, fmt.Errorf("Failed to generate secret key")
	}

	return sk, nil
}

func Base64BlsSk(sk *blst.SecretKey) (string, error) {
	skBytes := sk.Serialize()
	if skBytes == nil {
		return "", fmt.Errorf("Failed to base64 secret key")
	}

	s := base64.StdEncoding.EncodeToString(skBytes)

	return s, nil
}

func Unbase64BlsSk(b64sk string) (*blst.SecretKey, error) {
	skBytes, err := base64.StdEncoding.DecodeString(b64sk)
	if err != nil {
		return nil, err
	}

	sk := new(blst.SecretKey).Deserialize(skBytes)
	if sk == nil {
		return nil, fmt.Errorf("Failed to deserialize secret key from base64")
	}

	return sk, nil
}

// BLS Public Key operations
func GenBlsPk(sk *blst.SecretKey) (*PublicKey, error) {
	pk := new(PublicKey).From(sk)
	if pk == nil {
		return nil, fmt.Errorf("Failed to generate public key")
	}

	return pk, nil
}

func Base64BlsPk(pk *PublicKey) (string, error) {
	pkBytes := pk.Compress()
	if pkBytes == nil {
		return "", fmt.Errorf("Failed to compress pubkey")
	}

	s := base64.StdEncoding.EncodeToString(pkBytes)

	return s, nil
}

func Unbase64BlsPk(b64pk string) (*PublicKey, error) {
	pubkeyBytes, err := base64.StdEncoding.DecodeString(b64pk)
	if err != nil {
		return nil, err
	}

	pk := new(PublicKey).Uncompress(pubkeyBytes)
	if pk == nil {
		return nil, fmt.Errorf("Failed to uncompress pubkey from base64")
	}

	return pk, nil
}
