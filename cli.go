package main

import (
	"github.com/kagux/go-remote-cli/client"
	"github.com/kagux/go-remote-cli/server"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
	"path/filepath"
)

type Options struct {
	ServerOptions *server.Options
	ClientOptions *client.Options
	IsServer      bool
}

const (
	appVersion = "0.0.4"
)

var (
	app      = kingpin.New("remote_cli", "A remote command line app proxy")
	isServer = app.Flag("server", "Run in server mode").Short('s').Default("false").Bool()
	cmd      = app.Flag("cmd", "[Client mode] Command to run").Short('c').String()
	host     = app.Flag("host", "Host to bind to or to connect to").Short('h').Default("0.0.0.0").String()
	port     = app.Flag("port", "Port to bind to or to connect to").Short('p').Default("9201").Int()
)

func ParseCLI() *Options {
	exec_name := filepath.Base(os.Args[0])

	app.Version(appVersion)
	app.DefaultEnvars()
	// keep name dynamic to have env vars like MAHOUT_PORT=999
  app.Name = exec_name

	_, err := app.Parse(os.Args[1:])

	if err != nil {
		// use executable name as command and args as cmd args
		*cmd = exec_name + " " + strings.Join(os.Args[1:], " ")
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
		},
	}
}
