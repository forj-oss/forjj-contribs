package main

//go:generate go build -o $GOPATH/bin/forjj-genapp forjj-contribs/upstream/github/vendor/github.com/forj-oss/goforjj/genapp
//go:generate forjj-genapp github.yaml vendor/github.com/forj-oss/goforjj/genapp
