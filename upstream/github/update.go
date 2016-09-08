// This file has been created by "go generate" as initial code and HAS been updated. Do not remove it.

package main

import (
    "github.hpe.com/christophe-larsonneur/goforjj"
)

func (g *GitHubStruct)update_yaml_data(req *UpdateReq, ret *goforjj.PluginData) (Updated bool) {
    if g.github_source.Repos == nil {
        g.github_source.Repos = make(map[string]RepositoryStruct)
    }

    for name, repo := range req.ReposData {
        if g.AddRepo(name, repo) {
            ret.StatusAdd("New Repository '%s' added.", name)
        } else {
            r := g.github_source.Repos[name]
            if r.Update(repo) >0 {
                ret.StatusAdd("Repository '%s' updated.", name)
            }
        }

    }

    return
}
