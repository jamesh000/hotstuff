package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jamesh000/hotstuff/crypt"
	"github.com/jamesh000/hotstuff/net"
	"github.com/libp2p/go-libp2p/core/crypto"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// arguments
	//port := flag.Int("l", 0, "Port to listen on")
	destFlag := flag.String("d", "", "Destination address")
	rsaKeyFlag := flag.String("k", "rsa.key", "RSA key to use for libp2p")
	flag.Parse()

	var priv crypto.PrivKey
	if _, err := os.Stat(*rsaKeyFlag); errors.Is(err, os.ErrNotExist) {
		log.Printf("Generating RSA private key at %q", *rsaKeyFlag)
		priv, err := crypt.GenRSAKey()
		if err != nil {
			panic(err)
		}

		err = crypt.SaveRSAKey(*rsaKeyFlag, priv)
		if err != nil {
			panic(err)
		}
	} else {
		priv, err = crypt.LoadRSAKey(*rsaKeyFlag)
		if err != nil {
			panic(err)
		}
	}

	r, err := net.NewResonance(ctx, priv, 2, "synchronicity", "1.0")
	if err != nil {
		panic(err)
	}

	if *destFlag != "" {
		err = r.Connect(ctx, *destFlag)
		if err != nil {
			panic(err)
		}
	}

	err = r.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(r)

	go func() {
		for {
			msg, err := r.NextBroadcast(ctx, 0)
			if err != nil {
				panic(err)
			}

			fmt.Println(string(msg))
		}
	}()

	for {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := scanner.Text()

			r.SendBroadcast(ctx, 0, []byte(msg))

		}
	}

}
