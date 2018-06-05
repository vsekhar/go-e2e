This is an experimental app written in Go on the server and client.

# Prerequisites

 - [Go compiler](https://golang.org)
 - [Docker](https://docker.com)

# Server

From the repo root, build the server container and run it with host networking:

```bash
$ docker build -f server/Dockerfile -t server .
...
$ docker run --rm -p 8080:8080 server
```

Use CTRL-C to terminate the server.

# Web client

Navigate to `localhost:8080` in Chrome.

# Android

From the repo root, build the Android container (which will also build the app):

```bash
$ docker build -f android/Dockerfile -t android .
```

Run the container to copy the outputs. The following places the outputs where Android Studio would have:

```bash
$ OUTPUT_ABS_PATH=${PWD}/android/app/build/outputs
$ docker run --rm --user $UID:$(id -g) -v ${OUTPUT_ABS_PATH}:/outputs android
```

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
