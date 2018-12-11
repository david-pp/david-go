package main

import (
	"fmt"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var port = flag.String("port", "4729", "port")
var dir = flag.String("dir", "./", "work directory")
var daemon = flag.Bool("d", false, "run app as a daemon with -d=true or -d true.")

func main() {
	flag.Parse()

	if *daemon {
		cmd := exec.Command(os.Args[0])
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		fmt.Printf("%s [PID] %d running...\n", os.Args[0], cmd.Process.Pid)
		*daemon = false
		os.Exit(0)
	}

	h := http.FileServer(http.Dir(*dir))
	log.Fatal(http.ListenAndServe(":" + *port, h))
}