package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jamesh000/hotstuff/net"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// arguments
	// Network flags
	//port := flag.Int("l", 0, "Port to listen on")
	destFlag := flag.String("d", "", "Destination address")

	// Local config flag
	confFlag := flag.String("conf", "", "Config location")

	// New local config flags
	newConfFlag := flag.String("newconf", "", "file to generate new config in")
	idFlag := flag.Int("id", 0, "Id of this node")
	flag.Parse()

	if (*confFlag == "" && *newConfFlag == "") || (*confFlag != "" && *newConfFlag != "") {
		log.Println("Flag confusion, try again")
		return
	}

	var lcFile string

	if *newConfFlag != "" {
		err := NewLocalConf(*newConfFlag, *idFlag)
		if err != nil {
			panic(err)
		}

		log.Printf("Generated new config at %s", *newConfFlag)

		lcFile = *newConfFlag
	} else {
		lcFile = *confFlag
	}

	lc, err := LoadLocalConf(lcFile)
	if err != nil {
		panic(err)
	}
	log.Printf("Loaded config from %s\n", lcFile)

	r, err := net.NewResonance(ctx, lc.RsaKey, 2, "synchronicity", "1.0")
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
