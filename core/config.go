package core

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jamesh000/hotstuff/crypt"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

type replika struct {
	pid    peer.ID
	pubkey crypt.PublicKey
}

type Config struct {
	ReplikaCount int

	Pubkeys   []*crypt.PublicKey
	allowList map[peer.ID]struct{}
}

func (c *Config) InterceptSecured(nd network.Direction, pid peer.ID, ncm network.ConnMultiaddrs) bool {
	_, allow := c.allowList[pid]
	return allow
}

type replikaConf struct {
	number     int
	pid        string
	pubkeyPath string
}

func LoadConfig(path string) (*Config, error) {
	var conf []replikaConf

	confText, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(confText, &conf)
	if err != nil {
		return nil, err
	}

	fmt.Println(conf)

	config := new(Config)
	config.ReplikaCount = len(conf)
	config.Pubkeys = make([]*crypt.PublicKey, config.ReplikaCount)

	for _, r := range conf {
		if r.number > config.ReplikaCount {
			return nil, fmt.Errorf("Invalid replika number %d", r.number)
		}

		p, err := peer.Decode(r.pid)
		if err != nil {
			return nil, err
		}
		config.allowList[p] = struct{}{}

		config.Pubkeys[r.number], err = crypt.LoadBLSPubkey(r.pubkeyPath)
		if err != nil {
			return nil, err
		}
	}

	fmt.Println(config)

	return config, nil
}
