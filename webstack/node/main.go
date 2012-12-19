package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Expected mode of `deploy` or `server`")
		os.Exit(1)
	}
	switch args[0] {
	case "deploy":
		deploy(args[1:])
	case "server":
		server(args[1:])
	default:
		fmt.Println("Unknown mode `" + args[0] + "`")
		os.Exit(1)
	}
}