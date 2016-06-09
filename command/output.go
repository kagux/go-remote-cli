package command

import (
	"fmt"
)

type Output struct {
	Text string
	ExitStatus int
}

func NewErrorOutput(err error) *Output {
	msg := fmt.Sprintf("Error: %v\n", err)
	exitStatus := 1
	return &Output{
		Text: msg,
		ExitStatus: exitStatus,
	}
}

type OutputWriter struct {
	out chan *Output
}

func (w *OutputWriter) Write(p []byte) (n int, err error) {
	w.out <- &Output{ Text: string(p) }

	return len(p), nil
}
