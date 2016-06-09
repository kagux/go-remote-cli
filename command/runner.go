package command

import (
	"os/exec"
	"strings"
)

type Runner struct {
}

func NewRunner() *Runner {
	return &Runner{}
}

func (cr *Runner) Run(cmdStr string, out chan *Output) {
	cmdParts := strings.Fields(cmdStr)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	writer := &OutputWriter{out: out}
	cmd.Stdout = writer
	cmd.Stderr = writer

	if err := cmd.Start(); err != nil {
		out <- NewErrorOutput(err)
	}

	if err := cmd.Wait(); err != nil {
		out <- NewErrorOutput(err)
	}
}
