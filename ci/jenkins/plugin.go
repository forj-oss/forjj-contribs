package main

//go:generate go build -o $GOPATH/bin/forjj-genapp forjj-contribs/ci/jenkins/vendor/github.com/forj-oss/goforjj/genapp
//go:generate forjj-genapp jenkins.yaml vendor/github.com/forj-oss/goforjj/genapp
