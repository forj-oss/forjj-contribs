// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

// You can remove following comments.
// It has been designed fo you, to implement the core of your plugin task.
//
// You can use use it to write your own plugin handler for additional functionnality
// Like Index which currently return a basic code.

import (
    "net/http"
    "github.hpe.com/christophe-larsonneur/goforjj"
    "log"
    "fmt"
    "path"
    "os"
)

const github_file = "github.yaml"

// Do creating plugin task
// req_data contains the request data posted by forjj. Structure generated from 'github.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoCreate(w http.ResponseWriter, r *http.Request, req *CreateReq, ret *goforjj.PluginData) (httpCode int){

    gws := GitHubStruct{
        source_mount: req.ForjjSourceMount,
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

    // Build gws.github_source yaml structure.
    if err := gws.create_yaml_data(req) ; err != nil {
        ret.Errorf("%s", err)
        return
    }

    // A create won't be possible if repo requested already exist. The Update is the only possible option.
    // The list of repository found are listed and returned in the answer.
    if err := gws.repos_exists(ret) ; err != nil {
        ret.Errorf("%s\nUnable to create the github configuration when github already has repositories created. Use 'update' instead.", err)
        return 419
    }

    // A create won't be possible if source files already exist. The Update is the only possible option.
    log.Printf("Checking Infrastructure code existence.")
    source_path := path.Join(req.ForjjSourceMount, req.ForjjInstanceName)
    if _, err := os.Stat(source_path) ; err != nil {
        if err = os.MkdirAll(source_path, 0755) ; err != nil {
            ret.Errorf("Unable to create '%s'. %s", source_path, err)
        }
    }
    if _, err := os.Stat(path.Join(req.ForjjSourceMount, req.ForjjInstanceName, github_file)) ; err == nil {
        ret.Errorf("Unable to create the github configuration which already exist.\nUse 'update' to update it (or update %s), and 'maintain' to update your github service according to his configuration.", path.Join(req.ForjjInstanceName, github_file))
        return 419
    }

    ret.StatusAdd("Environment checked. Ready to be created.")

    // Save gws.github_source.
    if err := gws.save_yaml(path.Join(source_path, github_file)) ; err != nil {
        ret.Errorf("%s", err)
        return
    }
    log.Printf(ret.StatusAdd("Configuration saved in '%s'.", path.Join(req.ForjjInstanceName, github_file)))

    // Building final Post answer
    // We assume ssh is used and forjj can push with appropriate credential.
    infra_repo := gws.github_source.Repos[req.ForjjInfra]
    ret.Repos[req.ForjjInfra] = goforjj.PluginRepo{
        Name: infra_repo.Name,
        Exist: infra_repo.exist,
        Remotes: infra_repo.remotes,
        BranchConnect: infra_repo.branchConnect,
    }
    for k, v := range gws.github_source.Urls {
        ret.Services.Urls[k] = v
    }

    ret.CommitMessage = fmt.Sprintf("Create github configuration")
    ret.Files = append(ret.Files, path.Join(req.ForjjInstanceName, github_file))

    return
}

// Do updating plugin task
// req_data contains the request data posted by forjj. Structure generated from 'github.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoUpdate(w http.ResponseWriter, r *http.Request, req *UpdateReq, ret *goforjj.PluginData) (httpCode int) {

    gws := GitHubStruct{
        source_mount: req.ForjjSourceMount,
    }
    check := make(map[string]bool)

    if gws.verify_req_fails(ret, check) {
        return 422
    }

    // Read the github.yaml file.
    if err := gws.load_yaml(path.Join(req.ForjjSourceMount, github_file)) ; err != nil {
        ret.Errorf("Unable to update github instance '%s' source files. %s. Use 'create' to create it first.", req.ForjjInstanceName, err)
        return
    }

    // TODO: Update github source code
    /* if err := gws.update_yaml_data(req) ; err != nil {
        ret.Errorf("%s", err)
        return
    }

    // Save gws.github_source.
    if err := gws.save_yaml(path.Join(source_path, github_file)) ; err != nil {
        ret.Errorf("%s", err)
        return
    }
    log.Printf(ret.StatusAdd("Configuration saved in '%s'.", path.Join(req.ForjjInstanceName, github_file)))

    // Building final Post answer
    // We assume ssh is used and forjj can push with appropriate credential.
    infra_repo := gws.github_source.Repos[req.ForjjInfra]
    ret.Repos[req.ForjjInfra] = goforjj.PluginRepo{
        Name: infra_repo.Name,
        Exist: infra_repo.Exist,
        Remotes: infra_repo.Remotes,
        BranchConnect: infra_repo.BranchConnect,
    }
    for k, v := range gws.github_source.Urls {
        ret.Services.Urls[k] = v
    }

    ret.CommitMessage = fmt.Sprintf("Create github configuration")*/
    return
}

// Do maintaining plugin task
// req_data contains the request data posted by forjj. Structure generated from 'github.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoMaintain(w http.ResponseWriter, r *http.Request, req *MaintainReq, ret *goforjj.PluginData) (httpCode int) {

    gws := GitHubStruct{
        source_mount: req.ForjjSourceMount,
        workspace_mount: req.ForjjWorkspaceMount,
        token: req.GithubToken,
    }
    check := make(map[string]bool)
    check["token"] = true
    check["workspace"] = true

    if gws.verify_req_fails(ret, check) { // true => include workspace testing.
        return
    }

    // Read the github.yaml file.
    if err := gws.load_yaml(path.Join(req.ForjjSourceMount, req.ForjjInstanceName, github_file)) ; err != nil {
        ret.Errorf("%s", err)
        return
    }

    if gws.github_connect(gws.github_source.Urls["github-base-url"], ret) == nil {
        return
    }

    // ensure organization exist
    if ! gws.ensure_organization_exists(ret) {
        return
    }
    log.Printf(ret.StatusAdd("Organization maintained."))

    // loop on list of repos, and ensure they exist with minimal config and rights
    for name, repo_data := range  gws.github_source.Repos {
        if err := repo_data.ensure_exists(&gws, ret) ; err != nil {
           return
        }
        log.Printf(ret.StatusAdd("Repo maintained: %s", name))
    }
    return
}
