package command

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Runner struct {
}

func NewRunner() *Runner {
	return &Runner{}
}

func (cr *Runner) Run(cmdStr string, writer *OutputWriter) {
	if len(cmdStr) == 0 {
		writer.WriteError(errors.New("Received empty command"))
		return
	}

	fmt.Println("Command Received:", cmdStr)
	cmdParts := strings.Fields(cmdStr)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = writer
	cmd.Stderr = writer

	if err := cmd.Start(); err != nil {
		writer.WriteError(err)
	}

	if err := cmd.Wait(); err != nil {
		writer.WriteError(err)
	}
}
