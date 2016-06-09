package main

//go:generate go run $GOPATH/src/github.hpe.com/christophe-larsonneur/go-forjj/cmd/genflags/main.go github.yaml

import (
        "fmt"
//        "gopkg.in/alecthomas/kingpin.v2"
//        "github.com/alecthomas/kingpin"
//        "net/url"
//        "path"
//        "regexp"
)

func ( *githubApp)create(){
  fmt.Printf("create")
}

func ( *githubApp)update(){
  fmt.Printf("update")
}

func ( *githubApp)maintain(){
  fmt.Printf("maintain")
}
