package net

import (
	"context"
	"fmt"
	"log"

	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

const rendevousName string = "/synchronicity/1.0"

type Resonance struct {
	H   host.Host
	dht *kaddht.IpfsDHT

	ps           *pubsub.PubSub
	broadcast    []*pubsub.Topic
	broadcastSub []*pubsub.Subscription
}

func (r Resonance) String() string {
	return getHostAddress(r.H)
}

func NewResonance(ctx context.Context, id crypto.PrivKey, channelCount int, proto string, version string) (*Resonance, error) {
	r := new(Resonance)
	var err error
	r.H, r.dht, err = makeRoutedHost(ctx, 0, id)
	if err != nil {
		return nil, err
	}

	// Start pubsub (GossipSub)
	r.ps, err = pubsub.NewGossipSub(ctx, r.H)
	if err != nil {
		return nil, err
	}

	// make some broadcasting channels and their corresponding subscriptions
	r.broadcast = make([]*pubsub.Topic, channelCount)
	r.broadcastSub = make([]*pubsub.Subscription, channelCount)

	for i := 0; i < channelCount; i++ {
		// make broadcast i
		r.broadcast[i], err = r.ps.Join(fmt.Sprintf("/%s-bc-c%d/%s", proto, i, version))
		if err != nil {
			return nil, err
		}

		// make sub i
		r.broadcastSub[i], err = r.broadcast[i].Subscribe()
		if err != nil {
			return nil, err
		}
	}

	return r, err
}

func (r *Resonance) Bootstrap(ctx context.Context) error {
	err := r.dht.Bootstrap(ctx)
	if err != nil {
		return err
	}

	routingDiscovery := drouting.NewRoutingDiscovery(r.dht)
	dutil.Advertise(ctx, routingDiscovery, rendevousName)

	routingPeerC, err := routingDiscovery.FindPeers(ctx, rendevousName)
	if err != nil {
		panic(err)
	}

	for peer := range routingPeerC {
		if peer.ID == r.H.ID() {
			continue
		}

		err := r.H.Connect(ctx, peer)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Connected to peer", peer.ID)
	}

	return nil
}

func (r *Resonance) Connect(ctx context.Context, dest string) error {
	targetAddr, err := multiaddr.NewMultiaddr(dest)
	if err != nil {
		return err
	}
	targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		return err
	}
	err = r.H.Connect(ctx, *targetInfo)
	if err != nil {
		return err
	}
	fmt.Println("Connected to", targetInfo.ID)
	return nil
}

func (r *Resonance) NextBroadcast(ctx context.Context, c int) ([]byte, error) {
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		msg, err := r.broadcastSub[c].Next(ctx)
		if err != nil {
			return nil, err
		}

		if msg.ReceivedFrom == r.H.ID() {
			continue
		}

		return msg.Data, err
	}
}

func (r *Resonance) SendBroadcast(ctx context.Context, c int, msg []byte) error {
	err := r.broadcast[c].Publish(ctx, msg)

	return err
}
