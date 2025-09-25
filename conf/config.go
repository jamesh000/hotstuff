package conf

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bits-and-blooms/bitset"
	"github.com/jamesh000/hotstuff/crypt"
	"github.com/libp2p/go-libp2p/core/control"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type ReplikaConf struct {
	ReplikaCount uint

	Pubkeys   []*crypt.PublicKey
	allowList map[peer.ID]struct{}
}

func (c *ReplikaConf) InterceptSecured(nd network.Direction, pid peer.ID, ncm network.ConnMultiaddrs) bool {
	_, allow := c.allowList[pid]
	return allow
}

func (c *ReplikaConf) InterceptPeerDial(pid peer.ID) bool {
	_, allow := c.allowList[pid]
	return allow
}

// Don't really care about these other methods
func (c *ReplikaConf) InterceptAccept(network.ConnMultiaddrs) bool {
	return true
}

func (c *ReplikaConf) InterceptAddrDial(peer.ID, multiaddr.Multiaddr) bool {
	return true
}

func (c *ReplikaConf) InterceptUpgraded(network.Conn) (bool, control.DisconnectReason) {
	return true, 0
}

type replikaJson struct {
	Id        uint   `json:"id"`
	PeerId    string `json:"peerid"`
	PublicKey string `json:"publickey"`
}

func LoadReplikaConf(path string) (*ReplikaConf, error) {
	var conf []replikaJson

	confText, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(confText, &conf)
	if err != nil {
		return nil, err
	}

	config := new(ReplikaConf)
	config.ReplikaCount = uint(len(conf))
	config.Pubkeys = make([]*crypt.PublicKey, config.ReplikaCount)
	config.allowList = make(map[peer.ID]struct{})

	idsUsed := bitset.New(uint(config.ReplikaCount))

	for _, r := range conf {
		if r.Id > config.ReplikaCount {
			return nil, fmt.Errorf("Invalid replika number %d", r.Id)
		}

		if idsUsed.Test(r.Id) {
			return nil, fmt.Errorf("Replika ID %d appears twice", r.Id)
		}

		idsUsed.Set(r.Id)

		p, err := peer.Decode(r.PeerId)
		if err != nil {
			return nil, err
		}
		config.allowList[p] = struct{}{}

		config.Pubkeys[r.Id], err = crypt.Unbase64BlsPk(r.PublicKey)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func GenReplikaConf(rConfFile string, localConfFiles []string) error {
	newConf := make([]replikaJson, len(localConfFiles))
	for i, f := range localConfFiles {
		localConfBytes, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		var localConf lcText
		err = json.Unmarshal(localConfBytes, &localConf)
		if err != nil {
			return err
		}

		newConf[i].Id = localConf.Id
		newConf[i].PeerId = localConf.PeerId
		newConf[i].PublicKey = localConf.ConsensusPk
	}

	confBytes, err := json.MarshalIndent(newConf, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(rConfFile, confBytes, 0600)
	if err != nil {
		return err
	}

	return nil
}
