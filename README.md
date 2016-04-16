# Remote CLI

## Introduction
Simple client-server application to execute shell commands on remote machine.
It's meant to be used in the context of docker containers to seamlessly stub local executables.

For example, let's say you have a ruby web application that needs to execute java tool for some calculations.
Instead of having both web app runtime and java runtime in same container, we can replace executable with `remote_cli`
and have java in separate container.

## Usage

```
usage: remote_cli [<flags>]

A remote command line app proxy

Flags:
      --help            Show context-sensitive help (also try --help-long and --help-man).
  -s, --server          Run in server mode
  -c, --cmd=CMD         [Client mode] Command to run
  -h, --host="0.0.0.0"  Host to bind to or to connect to
  -p, --port=9201       Port to bind to or to connect to
      --version         Show application version.
```

In client mode you can either pass full command to `-c` flag, i.e. `remote_cli -c 'ls -la' -h server -p 9999`.
But a more convinient way is to rename `remote_cli` binary to executable you're stubbing. 
`remote_cli` will pick command name from it's file name and pass all arguments as is to remote server.
You can pass additional flags using environment variables as `[EXECUTABLE_NAME]_[FLAG]=value`.

### Simple example

Here's how you can replace `ls` command:

1. Create server image
  ```
    # server docker image
    FROM alpine

    ENV RC_VERSION 0.0.6
    RUN wget https://github.com/kagux/go-remote-cli/releases/download/${RC_VERSION}/linux-amd64-remote_cli.tar.bz2 \
        && tar -jxvf linux-amd64-remote_cli.tar.bz2 \
        && mv bin/linux/amd64/remote_cli /usr/local/bin/remote_cli

    EXPOSE 9021
    CMD ["remote_cli", "--server", "--host=0.0.0.0", "--port=9999"]
  ```

2. Run server container `docker run -d server`

3. Create client container
  ```
    # client docker image
    FROM alpine

    ENV RC_VERSION 0.0.6
    RUN wget https://github.com/kagux/go-remote-cli/releases/download/${RC_VERSION}/linux-amd64-remote_cli.tar.bz2 \
        && tar -jxvf linux-amd64-remote_cli.tar.bz2 \
        && mv bin/linux/amd64/remote_cli /bin/ls
  ```
4. Run client and list files on remote machine 
`docker run --rm -e LS_PORT=9999 -e LS_HOST=server --link server:server client ls`


## Release

Run `GITHUB_TOKEN=[your_token] make release`
