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
}

func New(opts *Options) *Client {
	return &Client{opts: opts}
}

func (c *Client) Run() error {
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		return err
	}
	fmt.Fprintf(conn, "%s\n", c.opts.Cmd)
	return nil
}
