package server

import (
	"encoding/gob"
	"fmt"
	"github.com/kagux/go-remote-cli/command"
	"net"
	"sync"
)

type RequestHandler struct {
	conn      net.Conn
	cmdRunner *command.Runner
	out       chan *command.Output
	waitGroup sync.WaitGroup
}

func NewRequestHandler(conn net.Conn) *RequestHandler {
	return &RequestHandler{
		conn:      conn,
		cmdRunner: command.NewRunner(),
	}
}

func (rh *RequestHandler) Handle() {
	rh.out = make(chan *command.Output)
	rh.waitGroup.Add(1)
	go rh.handleCommandOutput()
	rh.executeCommand()
	rh.waitGroup.Wait()
	rh.conn.Close()
	fmt.Println("*** Connection closed")
}

func (rh *RequestHandler) executeCommand() {
	dec := gob.NewDecoder(rh.conn)
	var r command.Request
	err := dec.Decode(&r)
	writer := command.NewOutputWriter(rh.out)
	if err != nil {
		writer.WriteError(err)
	}
	writer.Quiet = r.Quiet
	rh.cmdRunner.Run(r.NormalizedCommand(), writer)
	close(rh.out)
}

func (rh *RequestHandler) handleCommandOutput() {
	enc := gob.NewEncoder(rh.conn)
	for o := range rh.out {
		fmt.Print(o.Text)
		enc.Encode(o)
	}
	rh.waitGroup.Done()
}
