package server

import (
	"encoding/gob"
	"fmt"
	"github.com/kagux/go-remote-cli/command"
	"net"
	"io/ioutil"
	"io"
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
	eWriter := command.NewErrorWriter(rh.out)
	req, err := rh.readCommandRequest()
	if err != nil {
		eWriter.Write(err)
	}
	if err := rh.runCommand(req); err != nil {
		eWriter.Write(err)
	}
	close(rh.out)
}

func (rh *RequestHandler) runCommand(r command.Request) error {
	var w io.Writer
	if r.Quiet {
		w = ioutil.Discard
	} else {
		w = command.NewOutputWriter(rh.out)
	}
	return rh.cmdRunner.Run(r.NormalizedCommand(), w)
}

func (rh *RequestHandler) readCommandRequest() (r command.Request, err error) {
	dec := gob.NewDecoder(rh.conn)
	err = dec.Decode(&r)
	return
}

func (rh *RequestHandler) handleCommandOutput() {
	enc := gob.NewEncoder(rh.conn)
	for o := range rh.out {
		fmt.Print(o.Text)
		enc.Encode(o)
	}
	rh.waitGroup.Done()
}
