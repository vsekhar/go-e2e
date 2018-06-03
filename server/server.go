package main

import (
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
)

var port = flag.Int("port", 8080, "server port")
var path = flag.String("path", ".", "path to serve")

func main() {
	flag.Parse()
	fileserver := http.FileServer(http.Dir(*path))
	mime.AddExtensionType(".wasm", "application/wasm")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), fileserver))
}
