package server

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type CommandRunner struct {
}

func NewCommandRunner() *CommandRunner {
	return &CommandRunner{}
}

func (cr *CommandRunner) Run(cmdStr string, output io.Writer) {
	pReader, pWriter := io.Pipe()
	defer pWriter.Close()
	cmdParts := strings.Fields(cmdStr)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = pWriter
	cmd.Stderr = pWriter

	go io.Copy(output, pReader)

	if err := cmd.Start(); err != nil {
		fmt.Fprintln(pWriter, "Error starting command: ", err.Error())
	}

	if err := cmd.Wait(); err != nil {
		fmt.Fprintln(pWriter, "Error executing command: ", err.Error())
	}
}
