package client

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
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
	var wg sync.WaitGroup
	wg.Add(1)
	go c.printOutput(conn, &wg)
	fmt.Fprintf(conn, "%s\n", c.opts.Cmd)
	wg.Wait()

	return nil
}

func (c *Client) printOutput(conn net.Conn, wg *sync.WaitGroup) {
	io.Copy(os.Stdout, conn)
	wg.Done()
}
