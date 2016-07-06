package server

import (
	"encoding/gob"
	"fmt"
	"github.com/kagux/go-remote-cli/command"
	"net"
	"sync"
	"io/ioutil"
	"io"
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
	eWriter := command.NewErrorWriter(rh.out)
	if err != nil {
		eWriter.WriteError(err)
	}
	var oWriter io.Writer
	if r.Quiet {
		oWriter = ioutil.Discard
	} else {
		oWriter = command.NewOutputWriter(rh.out)
	}
	rh.cmdRunner.Run(r.NormalizedCommand(), oWriter, eWriter)
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
