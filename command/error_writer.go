package command

import (
	"fmt"
)

type ErrorWriter struct {
	out chan *Output
}

func NewErrorWriter(out chan *Output) *ErrorWriter {
	return &ErrorWriter{out: out}
}

func (w *ErrorWriter) WriteError(err error) {
	msg := fmt.Sprintf("Error: %v\n", err)
	exitStatus := 1
	w.out <- &Output{
		Text:       msg,
		ExitStatus: exitStatus,
	}
}
