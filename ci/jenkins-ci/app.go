// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
    "os/signal"
    "syscall"
    "net"
    "os"
    "log"
    "fmt"
    "net/http"
    "gopkg.in/alecthomas/kingpin.v2"
    "path"
)

func (a *jenkinsApp)start_server() {

    a.server_set()
    int_sig := make(chan os.Signal, 1)
    signal.Notify(int_sig, syscall.SIGINT, syscall.SIGTERM)

    ln, err := net.Listen("unix", a.socket)
    //ln, err := net.Listen("tcp", ":8081")
    if err != nil {
        log.Fatal("listen error:", err)
    }
    go func () {
        <-int_sig
        fmt.Printf("\nExiting and closing socket...\n")
        ln.Close()
        os.Exit(0)
    }()
    fmt.Printf("Starting http service on socket : %s\n", a.socket)

    router := NewRouter()

    srv := http.Server{Handler: router}
    log.Fatal(srv.Serve(ln))
}

func (a *jenkinsApp)server_set() {
    if _, err := os.Stat(*a.params.socket_path) ; err != nil {
        if os.IsNotExist(err) {
            os.Mkdir(*a.params.socket_path, 0755)
        } else {
            kingpin.FatalIfError(err, "Unable to create '%s'\n", *a.params.socket_path)
        }
    }
    a.socket = path.Join(*a.params.socket_path, *a.params.socket_file)
}

