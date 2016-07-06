package command

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
