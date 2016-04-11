package server

import (
	"strings"
	"os/exec"
	"fmt"
	"bufio"
)


type CommandRunner struct {
}


func (cr *CommandRunner) Run(cmdStr string, out chan string) error {
	cmdParts := strings.Fields(cmdStr)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			out <- fmt.Sprintf("%s\n", scanner.Text())
		}
	}()

	
	if err = cmd.Start(); err != nil {
		return err
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}
