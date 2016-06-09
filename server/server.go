package server

import (
	"fmt"
	"net"
)

type Server struct {
	opts      *Options
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
		go NewRequestHandler(conn).Handle()
	}
}
