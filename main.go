package main

import (
	"./client"
	"./server"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

func main() {
	isServer := kingpin.Flag("server", "Run in server mode").Short('s').Bool()
	cmdArgs := kingpin.Arg("cmdArgs", "In client mode pass command arguments").String()
	kingpin.Version("0.0.1")
	kingpin.Parse()

	var err error
	if *isServer {
		s := server.New(&server.Options{})
		err = s.Run()
	} else {
		c := client.New(&client.Options{
			Args: *cmdArgs,
		})
		err = c.Run()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
