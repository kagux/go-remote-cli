package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Server struct {
	opts      *Options
	cmdRunner *CommandRunner
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
		cmdRunner: NewCommandRunner(),
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
	output := io.MultiWriter(os.Stdout, conn)
	err = s.cmdRunner.Run(cmd, output)
	if err != nil {
		fmt.Fprintln(output, "Error executing command: ", err.Error())
	}
	// give client some time to read output
	time.Sleep(500 * time.Millisecond)
	conn.Close()
}
