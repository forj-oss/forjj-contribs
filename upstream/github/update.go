// This file has been created by "go generate" as initial code and HAS been updated. Do not remove it.

package main

import (
    "github.com/forj-oss/goforjj"
    "log"
)

func (g *GitHubStruct)update_yaml_data(req *UpdateReq, ret *goforjj.PluginData) (Updated bool) {
    if g.github_source.Repos == nil {
        g.github_source.Repos = make(map[string]RepositoryStruct)
    }

    log.Printf("Request has %d repository(ies)", len(req.ReposData))

    for name, repo := range req.ReposData {
        if g.AddRepo(name, repo) {
            Updated = true
            ret.StatusAdd("New Repository '%s' added.", name)
        } else {
            r := g.github_source.Repos[name]
            if r.Update(repo) >0 {
                Updated = true
                ret.StatusAdd("Repository '%s' updated.", name)
            }
        }

    }

    return
}

