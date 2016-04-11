package client

import (
	"fmt"
	"net"
	"bufio"
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
		// can have err and str at the same time
		if len(str) > 0 {
				fmt.Print(str)
		}
		if err == nil {
			continue
		}
		fmt.Println("Connection closed", err)
		c.wait <- true
		return
	}
}

func (c *Client) waitConn() {
	<- c.wait
}
