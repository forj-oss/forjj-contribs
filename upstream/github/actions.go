// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

// You can remove following comments.
// It has been designed fo you, to implement the core of your plugin task.
//
// You can use use it to write your own plugin handler for additional functionnality
// Like Index which currently return a basic code.

import (
//    "fmt"
//    "os"
    "net/http"
    "github.hpe.com/christophe-larsonneur/goforjj"
//    "github.com/parnurzeal/gorequest"
    "log"
)

// Do creating plugin task
// req_data contains the request data posted by forjj. Structure generated from 'github.yaml'.
// ret_data contains the response structure to return back to forjj.
//
func DoCreate(w http.ResponseWriter, r *http.Request, req *CreateReq, ret *goforjj.PluginData) {

    gws := GitHubStruct{
        source: req.ForjjSourceMount,
        token: req.GithubToken,
    }
    check := make(map[string]bool)
    check["token"] = true
    log.Printf("Checking parameters : %#v", gws)

    //ensure source path is writeable
    if gws.verify_req_fails(ret, check) {
        return
    }
    log.Printf("Checking github connection : %#v", gws)

    if gws.github_connect(req.GithubServer, ret) == nil {
        return
    }

    // Write the github.yaml source file.
}

// Do updating plugin task
// req_data contains the request data posted by forjj. Structure generated from 'github.yaml'.
// ret_data contains the response structure to return back to forjj.
//
func DoUpdate(w http.ResponseWriter, r *http.Request, req *UpdateReq, ret *goforjj.PluginData) {

    gws := GitHubStruct{
        source: req.ForjjSourceMount,
    }
    check := make(map[string]bool)

    if gws.verify_req_fails(ret, check) {
        return
    }
}

// Do maintaining plugin task
// req_data contains the request data posted by forjj. Structure generated from 'github.yaml'.
// ret_data contains the response structure to return back to forjj.
//
func DoMaintain(w http.ResponseWriter, r *http.Request, req *MaintainReq, ret *goforjj.PluginData) {

    gws := GitHubStruct{
        source: req.ForjjSourceMount,
        workspace: req.ForjjWorkspaceMount,
        token: req.GithubToken,
    }
    check := make(map[string]bool)
    check["token"] = true
    check["workspace"] = true

    if gws.verify_req_fails(ret, check) { // true => include workspace testing.
        return
    }

    // Read the github.yaml file.

    if gws.github_connect("", ret) == nil {
        return
    }
}
