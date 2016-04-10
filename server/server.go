package server

import (
	"bufio"
	"fmt"
	"net"
)

type Server struct {
}

type Options struct {
}

func New(opts *Options) *Server {
	return &Server{}
}

func (s *Server) Run() error {
	fmt.Println("Launching server...")
	ln, err := net.Listen("tcp", "localhost:8081")
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
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Print("Message Received:", string(message))
	conn.Close()
}
