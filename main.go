package main

import (
	"fmt"
	"github.com/kagux/go-remote-cli/client"
	"github.com/kagux/go-remote-cli/server"
	"os"
)

func main() {
	opts := ParseCLI(os.Args)
	var err error
	if opts.IsServer {
		err = server.New(opts.ServerOptions).Run()
	} else {
		err = client.New(opts.ClientOptions).Run()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "*** Error: %v\n", err)
		os.Exit(1)
	}
}
