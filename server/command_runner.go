package server

import (
	"io"
	"os/exec"
	"strings"
)

type CommandRunner struct {
}

func NewCommandRunner() *CommandRunner {
	return &CommandRunner{}
}

func (cr *CommandRunner) Run(cmdStr string, output io.Writer) error {
	cmdParts := strings.Fields(cmdStr)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = output
	cmd.Stderr = output

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
