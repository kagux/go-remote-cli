package server

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"github.com/kagux/go-remote-cli/command"
	"encoding/gob"
)

type Server struct {
	opts      *Options
	cmdRunner *command.Runner
}

type Options struct {
	Host string
	Port int
}

func (o *Options) Address() string {
	addr := fmt.Sprintf("%s:%d", o.Host, o.Port)
	return addr
}

func New(opts *Options) *Server {
	return &Server{
		opts:      opts,
		cmdRunner: command.NewRunner(),
	}
}

func (s *Server) Run() error {
	fmt.Println("Launching server...")
	ln, err := net.Listen("tcp", s.opts.Address())
	defer ln.Close()
	if err != nil {
		return err
	}
	fmt.Println("Listening...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go s.handleRequest(conn)
	}
}

func (s *Server) handleRequest(conn net.Conn) {
	cmd, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Print("Command Received:", string(cmd))
	out := make(chan *command.Output)
	go handleCommandOutput(conn, out)
	s.cmdRunner.Run(cmd, out)
	time.Sleep(500 * time.Millisecond)
	conn.Close()
	fmt.Println("Command executed")
}

func handleCommandOutput(conn net.Conn, out chan *command.Output) {
	enc := gob.NewEncoder(conn)
	for o := range out {
		fmt.Print(o.Text)
		enc.Encode(o)
	}
}
