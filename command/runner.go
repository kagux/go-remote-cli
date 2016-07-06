package command

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"io"
)

type Runner struct {
}

func NewRunner() *Runner {
	return &Runner{}
}

func (cr *Runner) Run(cmdStr string, oWriter io.Writer, eWriter *ErrorWriter) {
	if len(cmdStr) == 0 {
		eWriter.WriteError(errors.New("Received empty command"))
		return
	}

	fmt.Println("*** Command Received:", cmdStr)
	cmdParts := strings.Fields(cmdStr)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = oWriter
	cmd.Stderr = oWriter

	if err := cmd.Start(); err != nil {
		eWriter.WriteError(err)
	}

	if err := cmd.Wait(); err != nil {
		eWriter.WriteError(err)
	}
}
