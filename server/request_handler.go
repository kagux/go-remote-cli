package server

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"github.com/kagux/go-remote-cli/command"
	"net"
	"sync"
	"strings"
)

type RequestHandler struct {
	conn net.Conn
	cmdRunner  *command.Runner
	out chan *command.Output
	waitGroup sync.WaitGroup
}

func NewRequestHandler(conn net.Conn) *RequestHandler {
	return &RequestHandler{
		conn: conn,
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
	fmt.Println("Connection closed")
}

func (rh *RequestHandler) executeCommand() {
	cmd, err := bufio.NewReader(rh.conn).ReadString('\n')
	writer := command.NewOutputWriter(rh.out)
	if err != nil {
		writer.WriteError(err)
	}
	rh.cmdRunner.Run(strings.TrimSpace(cmd), writer)
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
