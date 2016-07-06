package command
import (
	"strings"
)

type Request struct {
	Cmd string
	Quiet bool
}

func (r *Request) NormalizedCommand() string {
	return strings.TrimSpace(r.Cmd)
}
