package main

import (
	"github.com/kagux/go-remote-cli/client"
	"github.com/kagux/go-remote-cli/server"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"path/filepath"
	"fmt"
)

type Options struct {
	ServerOptions *server.Options
	ClientOptions *client.Options
	IsServer      bool
}

const (
	appVersion = "0.0.6"
)

var (
	app      = kingpin.New("remote_cli", "A remote command line app proxy")
	isServer = app.Flag("server", "Run in server mode").Short('s').Default("false").Bool()
	host     = app.Flag("host", "Host to bind to or to connect to").Short('h').Default("0.0.0.0").String()
	port     = app.Flag("port", "Port to bind to or to connect to").Short('p').Default("9201").Int()
	cmd      = app.Flag("cmd", "[Client mode] Command to run").Short('c').String()
	quite    = app.Flag("quite", "[Client mode] Suppress command output").Short('q').Default("false").Bool()
)

func ParseCLI(args []string) *Options {
	exec_name := filepath.Base(args[0])
	cli_args := args[1:]

	app.Version(appVersion)
	app.DefaultEnvars()
	// keep name dynamic to have env vars like MAHOUT_PORT=999
  app.Name = exec_name

	_, err := app.Parse(cli_args)

	if err != nil {
		fmt.Println("*** Error parsing command: " + err.Error() + " ...passing all args to remote cli.")
		// use executable name as command and args as cmd args
		*cmd = exec_name + " " + strings.Join(cli_args, " ")
	}

	return &Options{
		IsServer: *isServer,
		ServerOptions: &server.Options{
			Host: *host,
			Port: *port,
		},
		ClientOptions: &client.Options{
			Cmd:  *cmd,
			Host: *host,
			Port: *port,
			Quite: *quite,
		},
	}
}
