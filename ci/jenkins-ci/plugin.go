package main

//go:generate go get github.hpe.com/christophe-larsonneur/goforjj gopkg.in/yaml.v2
//go:generate go build -o $GOPATH/bin/forjj-genapp github.hpe.com/christophe-larsonneur/goforjj/genapp
//go:generate forjj-genapp jenkins-ci.yaml
