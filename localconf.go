package main

import (
	"encoding/json"
	"os"

	"github.com/jamesh000/hotstuff/crypt"
	"github.com/libp2p/go-libp2p/core/crypto"
	blst "github.com/supranational/blst/bindings/go"
)

const rsaKeyFile string = "libp2p.key"
const blsSkFile string = "bls-sk.key"
const blsPkFile string = "bls-pk.key"
const confFile string = "config.json"

type LocalConf struct {
	Id           int
	RsaKey       crypto.PrivKey
	consensusKey *blst.SecretKey
}

type lcText struct {
	Id          int    `json:"id"`
	RsaKey      string `json:"rsakey"`
	ConsensusSk string `json:"consensuskey"`
	ConsensusPk string `json:"consensuspub"`
}

func LoadLocalConf(confFile string) (*LocalConf, error) {
	var text lcText

	confBytes, err := os.ReadFile(confFile)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(confBytes, &text); err != nil {
		return nil, err
	}

	lc := new(LocalConf)

	lc.Id = text.Id

	// load rsa key
	lc.RsaKey, err = crypt.Unbase64RSAKey(text.RsaKey)
	if err != nil {
		return nil, err
	}

	// load bls key
	lc.consensusKey, err = crypt.Unbase64BlsSk(text.ConsensusSk)
	if err != nil {
		return nil, err
	}

	return lc, nil
}

func NewLocalConf(confFile string, id int) error {
	var text lcText

	// Make local text
	text.Id = id

	// RSA key
	priv, err := crypt.GenRSAKey()
	if err != nil {
		return err
	}

	text.RsaKey, err = crypt.Base64RSAKey(priv)
	if err != nil {
		return err
	}

	// BLS secret key
	sk, err := crypt.GenBlsSk()
	if err != nil {
		return err
	}

	text.ConsensusSk, err = crypt.Base64BlsSk(sk)
	if err != nil {
		return err
	}

	// BLS public key
	pk, err := crypt.GenBlsPk(sk)
	if err != nil {
		return err
	}

	text.ConsensusPk, err = crypt.Base64BlsPk(pk)
	if err != nil {
		return err
	}

	textBytes, err := json.MarshalIndent(text, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(confFile, textBytes, 0600)
	if err != nil {
		return err
	}

	return nil
}
