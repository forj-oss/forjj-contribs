// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
    "gopkg.in/alecthomas/kingpin.v2"
    "os"
)

var cliApp jenkinsApp

func main() {
    cliApp.init()

    switch kingpin.MustParse(cliApp.App.Parse(os.Args[1:])) {
    case "service start":
        cliApp.start_server()
    default:
        kingpin.Usage()
    }
}
