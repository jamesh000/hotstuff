package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jamesh000/hotstuff/conf"
	"github.com/jamesh000/hotstuff/net"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// arguments
	// Network flags
	//port := flag.Uint("l", 0, "Port to listen on")
	destFlag := flag.String("d", "", "Destination address")

	// Local config flag
	confFlag := flag.String("conf", "", "Config location")

	// New local config flags
	newConfFlag := flag.String("newconf", "", "file to generate new config in")
	idFlag := flag.Uint("id", 0, "Id of this node")

	// Replika config flag
	rConfFlag := flag.String("rconf", "", "Replika config location")

	// Replika config generation flags
	newRcFlag := flag.String("newrconf", "", "New replika config location")
	newRcInputsFlag := flag.String("rcgin", "", "New replika config input local confs")

	flag.Parse()

	// Generation of replika configs
	if *newRcFlag != "" {
		log.Printf("Generating a new replika config at %s and exiting\n", *newRcFlag)

		if *newRcInputsFlag == "" {
			log.Printf("Please provide inputs to generate based off of\n")
			return
		}

		inputs := strings.Split(*newRcInputsFlag, ",")

		err := conf.GenReplikaConf(*newRcFlag, inputs)
		if err != nil {
			panic(err)
		}

		log.Println("Successfully created new config. Check it an then run with -rconf")
		return
	}

	// Generation of local configs
	if *newConfFlag != "" {
		err := conf.NewLocalConf(*newConfFlag, *idFlag)
		if err != nil {
			panic(err)
		}

		log.Printf("Generated new local config at %s", *newConfFlag)
		return
	}

	// Load local config
	lc, err := conf.LoadLocalConf(*confFlag)
	if err != nil {
		panic(err)
	}
	log.Printf("Loaded local config from %s\n", *confFlag)

	// Load replika config
	rc, err := conf.LoadReplikaConf(*rConfFlag)
	if err != nil {
		panic(err)
	}
	log.Printf("Loaded replika config %s\n", *rConfFlag)

	r, err := net.NewResonance(ctx, lc.RsaKey, rc, 2, "synchronicity", "1.0")
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
