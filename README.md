This is an experimental app written in Go on the server and client.

# Prerequisites

 - [Go compiler](https://golang.org)
 - [Docker](https://docker.com)

# Server

From the repo root, build the server container and run it with host networking:

```bash
$ docker build -f server/DockerFile -t server .
...
$ docker run --rm --network=host server
```

Use CTRL-C to terminate the server.

# Web client

Navigate to `localhost:8080` in Chrome.

# Android

TBD

# iOS

TBD

# IDE

For IDE auto-completion, install the go-wasm compiler to your system. You will need a regular [Go compiler](https://golang.org/) installed on your system and the go-wasm compiler will live alongside it.

```bash
$ git clone --branch wasm-wip https://github.com/neelance/go.git $HOME/go-wasm
$ cd $HOME/go-wasm/src && ./all.bash
```

To use the WASM-capable compiler for the duration of a terminal session:

```bash
$ GOROOT="$HOME/go-wasm"
$ alias go="$HOME/go-wasm/bin/go"
```

Configure VS Code by adding the following to your workspace settings (you'll need to expand out $HOME yourself):

```json
"go.goroot": "$HOME/go-wasm"
```

# Resources

 - https://github.com/neelance/go
 - https://blog.lazyhacker.com/2018/05/webassembly-wasm-with-go.html
