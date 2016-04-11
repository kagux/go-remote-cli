package server

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Server struct {
	opts *Options
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
	return &Server{ opts: opts }
}

func (s *Server) Run() error {
	fmt.Println("Launching server...")
	ln, err := net.Listen("tcp", s.opts.Address())
	defer ln.Close()
	if err != nil {
		return err
	}
	cmdRunner := &CommandRunner{}

	fmt.Println("Listening...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go s.handleRequest(conn, cmdRunner)
	}
}

func (s *Server) handleRequest(conn net.Conn, cmdRunner *CommandRunner) {
	cmd, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Print("Command Received:", string(cmd))
	out := make(chan string)
	go func() {
		for s := range out {
			fmt.Fprintf(conn, "%s", s)
			fmt.Print(s)
		}
	}()
	err = cmdRunner.Run(cmd, out)
	// give client some time to read output
	time.Sleep(500 * time.Millisecond)
	conn.Close()
}
