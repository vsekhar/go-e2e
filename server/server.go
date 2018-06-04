package main

import (
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
)

var port = flag.Int("port", 0, "server port (default: PORT environment variable)")
var path = flag.String("path", ".", "path to serve")

func main() {
	flag.Parse()
	if *port == 0 {
		portString, ok := os.LookupEnv("PORT")
		if !ok {
			panic("port must be specified using -port or PORT envvar")
		}
		var err error
		*port, err = strconv.Atoi(portString)
		if err != nil {
			panic("bad port envvar: " + portString)
		}
	}

	fileserver := http.FileServer(http.Dir(*path))
	mime.AddExtensionType(".wasm", "application/wasm")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), fileserver))
}
