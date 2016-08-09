// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

// You can remove following comments.
// It has been designed fo you, to implement the core of your plugin task.
//
// You can use use it to write your own plugin handler for additional functionnality
// Like Index which currently return a basic code.

import (
    "fmt"
    "os"
    "net/http"
    "github.hpe.com/christophe-larsonneur/goforjj"
)

// Do creating plugin task
// req_data contains the request data posted by forjj. Structure generated from 'jenkins.yaml'.
// ret_data contains the response structure to return back to forjj.
//
func DoCreate(w http.ResponseWriter, r *http.Request, req_data *CreateReq, ret_data *goforjj.PluginData) {

    // This is where you shoud write your Update code. Following line is for Demo only.
    fmt.Fprintf(os.Stdout,"%#v\n", req_data)

}

// Do updating plugin task
// req_data contains the request data posted by forjj. Structure generated from 'jenkins.yaml'.
// ret_data contains the response structure to return back to forjj.
//
func DoUpdate(w http.ResponseWriter, r *http.Request, req_data *UpdateReq, ret_data *goforjj.PluginData) {

    // This is where you shoud write your create code. Following line is for Demo only.
    fmt.Fprintf(os.Stdout,"%#v\n", req_data)

}

// Do maintaining plugin task
// req_data contains the request data posted by forjj. Structure generated from 'jenkins.yaml'.
// ret_data contains the response structure to return back to forjj.
//
func DoMaintain(w http.ResponseWriter, r *http.Request, req_data *MaintainReq, ret_data *goforjj.PluginData) {

    // This is where you shoud write your Update code. Following line is for Demo only.
    fmt.Fprintf(os.Stdout,"%#v\n", req_data)

}
