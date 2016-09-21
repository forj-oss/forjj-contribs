// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.
package main

import (
    "net/http"
    "github.hpe.com/christophe-larsonneur/goforjj"
)

// Do creating plugin task
// req_data contains the request data posted by forjj. Structure generated from 'jenkins.yaml'.
// ret_data contains the response structure to return back to forjj.
//
func DoCreate(w http.ResponseWriter, r *http.Request, req *CreateReq, ret *goforjj.PluginData) (httpCode int) {
    var p *JenkinsPlugin

    if pr, code := req.check_source_existence(ret) ; pr == nil {
        return code
    } else {
        p = pr
    }

    if ! p.initialize_from(req, ret) {
        return
    }

    if ! p.create_jenkins_sources(req.Args.ForjjInstanceName, ret) {
        return
    }

    if ! p.save_yaml(ret) {
        return
    }
    return
}

// Do updating plugin task
// req_data contains the request data posted by forjj. Structure generated from 'jenkins.yaml'.
// ret_data contains the response structure to return back to forjj.
//
func DoUpdate(w http.ResponseWriter, r *http.Request, req *UpdateReq, ret *goforjj.PluginData) (httpCode int) {
    var p *JenkinsPlugin

    if pr, ok := req.check_source_existence(ret) ; !ok {
        return
    } else {
        p = pr
    }

    if ! p.load_yaml(p.yaml.Settings.InstanceName, ret) {
        return
    }

    if ! p.update_from(req, ret) {
        return
    }

    if ! p.update_jenkins_sources(ret) {
        return
    }

    if ! p.save_yaml(ret) {
        return
    }
    return
}

// Do maintaining plugin task
// req_data contains the request data posted by forjj. Structure generated from 'jenkins.yaml'.
// ret_data contains the response structure to return back to forjj.
//
func DoMaintain(w http.ResponseWriter, r *http.Request, req *MaintainReq, ret *goforjj.PluginData) (httpCode int) {
    var p *JenkinsPlugin

    if pr, ok := req.check_source_existence(ret) ; !ok {
        return
    } else {
        p = pr
    }

    // loop on list of jenkins instances defined by a collection of */jenkins.yaml
    if ! p.InstantiateAll(ret) {
        return
    }
    return
}
