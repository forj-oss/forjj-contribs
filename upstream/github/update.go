// This file has been created by "go generate" as initial code and HAS been updated. Do not remove it.

package main

import (
	"fmt"
	"github.com/forj-oss/goforjj"
	"log"
)

func (g *GitHubStruct) update_yaml_data(req *UpdateReq, ret *goforjj.PluginData) (bool, error) {
	if g.github_source.Urls == nil {
		return false, fmt.Errorf("Internal Error. Urls was not set.")
	}

	if g.github_source.Repos == nil {
		g.github_source.Repos = make(map[string]RepositoryStruct)
	}

	// In update, we simply rebuild Users and Team from Forjfile.
	// No need to keep track of removed one
	g.github_source.Users = make(map[string]string)
	g.github_source.Groups = make(map[string]TeamStruct)

	if g.app.Repos_disabled == "true" {
		log.Print("Repos_disabled is true. forjj_github won't manage repositories except the infra one.")
		g.github_source.NoRepos = true
	} else {
		// Updating all from Forjfile repos
		g.github_source.NoRepos = false
		for name, repo := range req.Objects.Repo {
			if !repo.IsValid(name, ret) {
				continue
			}

			g.SetRepo(&repo, (name == g.app.ForjjInfra))
		}

		// Disabling missing one
		for name, repo := range g.github_source.Repos {
			if err := repo.IsValid(name); err != nil {
				delete(g.github_source.Repos, name)
				ret.StatusAdd("Warning!!! Invalid repository '%s' found in github.yaml. Removed.")
				continue
			}
			if _, found := req.Objects.Repo[name]; !found && !repo.Disabled {
				repo.Disabled = true
				g.github_source.Repos[name] = repo
				ret.StatusAdd("Disabling repository '%s'", name)
			}
		}
	}

	log.Printf("Github manage %d repository(ies)", len(g.github_source.Repos))

	if g.app.Teams_disabled == "true" {
		log.Print("Teams_disabled is true. forjj_github won't manage Organization Users.")
		g.github_source.NoTeams = true
	} else {
		g.github_source.NoTeams = false
		for name, details := range req.Objects.User {
			g.AddUser(name, &details)
		}
	}

	log.Printf("Github manage %d user(s) at Organization level.", len(g.github_source.Users))

	if g.github_source.NoTeams {
		log.Print("Teams_disabled is true. forjj_github won't manage Organization Groups.")
	} else {
		for name, details := range req.Objects.Group {
			g.AddGroup(name, &details)
		}
	}

	log.Printf("Github manage %d group(s) at Organization level.", len(g.github_source.Groups))

	return true, nil
}

// SetRepo Add a new repository to be managed by github plugin.
func (g *GitHubStruct) SetRepo(repo *RepoInstanceStruct, is_infra bool) {
	upstream := g.DefineRepoUrls(repo.Name)

	// found or not, I need to set it.
	r := RepositoryStruct{}
	r.set(repo,
		map[string]goforjj.PluginRepoRemoteUrl{"origin": upstream},
		map[string]string{"master": "origin/master"},
		is_infra)
	g.github_source.Repos[repo.Name] = r
}

// SaveMaintainOptions Function which adds maintain options as part of the plugin answer in create/update phase.
// forjj won't add any driver name because 'maintain' phase read the list of drivers to use from forjj-maintain.yml
// So --git-us is not available for forjj maintain.
func (r *UpdateArgReq) SaveMaintainOptions(ret *goforjj.PluginData) {
	if ret.Options == nil {
		ret.Options = make(map[string]goforjj.PluginOption)
	}
}

func addMaintainOptionValue(options map[string]goforjj.PluginOption, option, value, defaultv, help string) goforjj.PluginOption {
	opt, ok := options[option]
	if ok && value != "" {
		opt.Value = value
		return opt
	}
	if !ok {
		opt = goforjj.PluginOption{Help: help}
		if value == "" {
			opt.Value = defaultv
		} else {
			opt.Value = value
		}
	}
	return opt
}
