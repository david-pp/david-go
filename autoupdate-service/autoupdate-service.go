package main

import (
	// "fmt"
	"flag"
	"log"
	"net/http"
)

var port = flag.String("port", "4729", "port")
var dir = flag.String("dir", "./", "work directory")

func main() {
	flag.Parse()

	h := http.FileServer(http.Dir(*dir))
	log.Fatal(http.ListenAndServe(":" + *port, h))
}