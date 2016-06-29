package client

import (
	"encoding/gob"
	"fmt"
	"github.com/kagux/go-remote-cli/command"
	"io"
	"net"
	"os"
	"sync"
)

type Client struct {
	opts *Options
}

type Options struct {
	Cmd   string
	Host  string
	Port  int
	Quiet bool
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
	fmt.Printf("*** Connection to remote cli at %s established\n", c.opts.Address())
	var wg sync.WaitGroup
	wg.Add(1)
	go c.handleOutput(conn, &wg)
	fmt.Printf("*** Executing command `%s`\n", c.opts.Cmd)
	fmt.Fprintf(conn, "%s\n", c.opts.Cmd)
	wg.Wait()
	fmt.Println("*** Command successfully executed, closing connection")

	return nil
}

func (c *Client) handleOutput(conn net.Conn, wg *sync.WaitGroup) {
	dec := gob.NewDecoder(conn)
	var o command.Output
	for {
		err := dec.Decode(&o)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("*** Decoding error: %v\n", err)
			os.Exit(1)
		}
		if !c.opts.Quiet {
			fmt.Print(o.Text)
		}
		if o.ExitStatus > 0 {
			os.Exit(o.ExitStatus)
		}
	}
	wg.Done()
}
