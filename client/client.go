package client

import (
	"fmt"
	"net"
	"bufio"
	"io"
)

type Client struct {
	opts *Options
	wait chan bool
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
	return &Client{
		opts: opts,
		wait: make(chan bool, 1),
	}
}

func (c *Client) Run() error {
	conn, err := net.Dial("tcp", c.opts.Address())
	if err != nil {
		return err
	}
	go c.printOutput(conn)
	fmt.Fprintf(conn, "%s\n", c.opts.Cmd)
	c.waitConn()

	return nil
}

func (c *Client) printOutput(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err == nil {
				fmt.Print(str)
				continue
		}
		if err != io.EOF {
			panic("error while reading from connection")
		}
		c.wait <- true
		return
	}
}

func (c *Client) waitConn() {
	<- c.wait
}
