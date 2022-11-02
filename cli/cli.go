package cli

import (
	"flag"
	"fmt"
	"github.com/78planet/nomadcoin/explorer"
	"github.com/78planet/nomadcoin/rest"
	"os"
)

func usage() {
	fmt.Printf("hello this is nomadcoin\n\n")
	fmt.Printf("input option\n\n")
	fmt.Printf("port: \n")
	fmt.Printf("mode: html or rest")

	os.Exit(0)
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Sets the port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}
}
