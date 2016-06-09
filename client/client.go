package client

import (
	"os"
	"io"
	"fmt"
	"net"
	"sync"
	"github.com/kagux/go-remote-cli/command"
	"encoding/gob"
)

type Client struct {
	opts *Options
}

type Options struct {
	Cmd  string
	Host string
	Port int
}

func (o *Options) Address() string {
	addr := fmt.Sprintf("%s:%d", o.Host, o.Port)
	return addr
}

func New(opts *Options) *Client {
	return &Client{
		opts: opts,
	}
}

func (c *Client) Run() error {
	conn, err := net.Dial("tcp", c.opts.Address())
	if err != nil {
		return err
	}
	defer conn.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	go handleOutput(conn, &wg)
	fmt.Fprintf(conn, "%s\n", c.opts.Cmd)
	wg.Wait()

	return nil
}

func handleOutput(conn net.Conn, wg *sync.WaitGroup) {
	dec := gob.NewDecoder(conn)
	var o command.Output
	for {
		err := dec.Decode(&o)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Decoding error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(o.Text)
		if o.ExitStatus > 0 {
			os.Exit(o.ExitStatus)
		}
	}
	wg.Done()
}
