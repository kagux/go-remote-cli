package client

import (
	"fmt"
	"net"
)

type Client struct {
	opts *Options
}

type Options struct {
	Cmd string
	Host string
	Port int
}

func (o *Options) Address() string {
	addr := fmt.Sprintf("%s:%d", o.Host, o.Port)
	return addr
}

func New(opts *Options) *Client {
	return &Client{opts: opts}
}

func (c *Client) Run() error {
	conn, err := net.Dial("tcp", c.opts.Address())
	if err != nil {
		return err
	}
	fmt.Fprintf(conn, "%s\n", c.opts.Cmd)
	return nil
}
