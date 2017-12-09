package main

import (
	"fmt"
	"github.com/forj-oss/goforjj"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func (g *GitHubStruct) create_yaml_data(req *CreateReq, ret *goforjj.PluginData) error {
	// Write the github.yaml source file.
	if g.github_source.Urls == nil {
		return fmt.Errorf("Internal Error. Urls was not set.")
	}

	req.InitOrganization(g)

	g.github_source.Repos = make(map[string]RepositoryStruct)
	g.github_source.Users = make(map[string]string)
	g.github_source.Groups = make(map[string]TeamStruct)

	g.github_source.NoRepos = (g.app.ReposDisabled == "true")
	if g.github_source.NoRepos {
		log.Print("Repositories_disabled is true. forjj_github won't manage repositories except the infra repository.")
	}

	// Add all repos
	for name, repo := range req.Objects.Repo {
		is_infra := (name == g.app.ForjjInfra)
		if g.github_source.NoRepos && ! is_infra {
			continue
		}
		if !repo.IsValid(name, ret) {
			ret.StatusAdd("Warning!!! Invalid repository '%s' requested. Ignored.")
			continue
		}
		g.SetRepo(&repo, is_infra)
		g.SetHooks(&repo, req.Objects.Webhooks)
	}
	g.SetOrgHooks(req.Objects.Webhooks)


	log.Printf("forjj-github manages %d repository(ies).", len(g.github_source.Repos))

	g.github_source.NoTeams = (g.app.TeamsDisabled == "true")
	if g.github_source.NoTeams {
		log.Print("Users_disabled is true. forjj_github won't manage Organization teams (Users/groups).")
	} else {
		for name, details := range req.Objects.User {
			g.AddUser(name, &details)
		}
	}

	log.Printf("forjj-github manages %d user(s) at Organization level.", len(g.github_source.Users))

	if !g.github_source.NoTeams {
		for name, details := range req.Objects.Group {
			g.AddGroup(name, &details)
		}
	}

	log.Printf("forjj-github manages %d group(s) at Organization level.", len(g.github_source.Groups))

	return nil
}

func (g *GitHubStruct) DefineRepoUrls(name string) (upstream goforjj.PluginRepoRemoteUrl) {
	upstream = goforjj.PluginRepoRemoteUrl{
		Ssh: "git@" + g.Client.BaseURL.Host + ":" + g.github_source.Organization + "/" + name + ".git",
		Url: g.github_source.Urls["github-url"] + "/" + g.github_source.Organization + "/" + name,
	}
	return
}

// AddUser Add a new repository to be managed by github plugin.
func (g *GitHubStruct) AddUser(name string, UserDet *UserInstanceStruct) bool {
	if _, found := g.github_source.Users[name]; !found {
		g.github_source.Users[name] = UserDet.Role
		return true // New added
	}
	return false
}

// AddGroup Add a new repository to be managed by github plugin.
func (g *GitHubStruct) AddGroup(name string, GroupDet *GroupInstanceStruct) bool {
	if _, found := g.github_source.Groups[name]; !found {
		g.github_source.Groups[name] = TeamStruct{Role: GroupDet.Role, Users: GroupDet.Members}
		return true // New added
	}
	return false
}

func (g *GitHubStruct) save_yaml(file string) (Updated bool, _ error) {

	d, err := yaml.Marshal(&g.github_source)
	if err != nil {
		return false, fmt.Errorf("Unable to encode github data in yaml. %s", err)
	}

	if d_before, err := ioutil.ReadFile(file); err != nil {
		Updated = true
	} else {
		Updated = (string(d) != string(d_before))
	}

	if !Updated {
		return
	}
	if err := ioutil.WriteFile(file, d, 0644); err != nil {
		return false, fmt.Errorf("Unable to save '%s'. %s", file, err)
	}
	return
}

func (g *GitHubStruct) load_yaml(file string) error {
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Unable to load '%s'. %s", file, err)
	}

	err = yaml.Unmarshal(d, &g.github_source)
	if err != nil {
		return fmt.Errorf("Unable to decode github data in yaml. %s", err)
	}
	return nil
}

func (r *CreateArgReq) SaveMaintainOptions(ret *goforjj.PluginData) {
	if ret.Options == nil {
		ret.Options = make(map[string]goforjj.PluginOption)
	}
}
