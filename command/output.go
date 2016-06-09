package command

import (
	"fmt"
)

type Output struct {
	Text       string
	ExitStatus int
}

type OutputWriter struct {
	out chan *Output
}

func NewOutputWriter(out chan *Output) *OutputWriter {
	return &OutputWriter{out: out}
}

func (w *OutputWriter) Write(p []byte) (n int, err error) {
	w.out <- &Output{Text: string(p)}

	return len(p), nil
}

func (w *OutputWriter) WriteError(err error) {
	msg := fmt.Sprintf("Error: %v\n", err)
	exitStatus := 1
	w.out <- &Output{
		Text:       msg,
		ExitStatus: exitStatus,
	}
}
