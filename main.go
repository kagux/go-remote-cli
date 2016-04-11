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
	cmd := kingpin.Flag("cmd", "[Client mode] Command to run").Short('c').String()
	host := kingpin.Flag("host", "Host to bind to or to connect to").Short('h').Default("0.0.0.0").String()
	port := kingpin.Flag("port", "Port to bind to or to connect to").Short('p').Default("9201").Int()
	kingpin.Version("0.0.1")
	kingpin.Parse()

	var err error
	if *isServer {
		s := server.New(&server.Options{
			Host: *host,
			Port: *port,
		})
		err = s.Run()
	} else {
		c := client.New(&client.Options{
			Cmd: *cmd,
			Host: *host,
			Port: *port,
		})
		err = c.Run()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
