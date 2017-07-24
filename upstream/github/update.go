// This file has been created by "go generate" as initial code and HAS been updated. Do not remove it.

package main

import (
	"github.com/forj-oss/goforjj"
	"log"
)

func (g *GitHubStruct) update_yaml_data(req *UpdateReq, ret *goforjj.PluginData) bool {
	if g.github_source.Repos == nil {
		g.github_source.Repos = make(map[string]RepositoryStruct)
	}

	log.Printf("Request has %d repository(ies)", len(req.Objects.Repo))

	for _, repo := range req.Objects.Repo {
		Updated, err_msg, mess := repo.DoUpdateIn(g)
		if Updated {
			ret.StatusAdd(mess)
		} else {
			ret.ErrorMessage = err_msg
		}
	}

	return true
}

// Function which adds maintain options as part of the plugin answer in create/update phase.
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
