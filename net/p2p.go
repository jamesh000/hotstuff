package net

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/multiformats/go-multiaddr"
)

func makeRoutedHost(ctx context.Context, port int, priv crypto.PrivKey) (host.Host, *kaddht.IpfsDHT, error) {
	// Create the address of the node
	nodeAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	// Set up routing
	var dht *kaddht.IpfsDHT
	newDHT := func(h host.Host) (routing.PeerRouting, error) {
		dht, err = kaddht.New(ctx, h, kaddht.Mode(kaddht.ModeServer))
		return dht, err
	}

	// Create node
	h, err := libp2p.New(
		libp2p.ListenAddrs(nodeAddr),
		libp2p.Identity(priv),
		libp2p.Routing(newDHT),
	)

	return h, dht, err
}

func getHostAddress(ha host.Host) string {
	// Build host multiaddress
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", ha.ID()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := ha.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}

// Bin of random stuff
/*
func subHandler(ctx context.Context, sub *pubsub.Subscription) {
	defer sub.Cancel()
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Println("Failed to read next message on topic", sub.Topic(), "due to", err)
			continue
		}

		fmt.Println(msg)
	}
}

func publishLoop(ctx context.Context, topic *pubsub.Topic, donec chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		err := topic.Publish(ctx, []byte(msg))
		if err != nil {
			fmt.Println("Failed to publish message on topic", topic.String(), "due to", err)
		}
	}
	donec <- struct{}{}
}

func (r *Resonance) SendTo(ctx context.Context, pid peer.ID, msg []byte) error {
	s, err := r.H.NewStream(ctx, pid, protocol.ID(directName))

	if err != nil {
		return err
	}

	_, err = s.Write(msg)

	return err
}

func (r *Resonance) handleDirect(s network.Stream) {
	msg := make([]byte, 1024)
	fmt.Println("Got new stream")

	b := bufio.NewReader(s)

	msgSize, err := b.Read(msg)
	if err != nil {
		panic(err)
	}

	fmt.Println(msg)

	r.Direct <- msg[:msgSize]
}
*/
