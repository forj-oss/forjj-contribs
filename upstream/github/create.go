package main

import (
    "gopkg.in/yaml.v2"
    "fmt"
    "io/ioutil"
    "github.com/forj-oss/goforjj"
)

func (g *GitHubStruct)create_yaml_data(req *CreateReq) error {
    // Write the github.yaml source file.
    g.github_source.Urls = make(map[string]string)
    g.github_source.Urls["github-base-url"] = g.Client.BaseURL.String()

    req.InitOrganization(g)

    if g.github_source.Repos == nil {
        g.github_source.Repos = make(map[string]RepositoryStruct)
    }

    for name, repo := range req.Objects.Repo {
        g.AddRepo(name, &repo )
    }

    // TODO: Be able to add several repos thanks to the request structure.
    return nil
}

// Add a new repository to be managed by github plugin.
func (g *GitHubStruct)AddRepo(name string, repo *RepoInstanceStruct) bool{
    upstream := "git@" + g.Client.BaseURL.Host + ":" + g.github_source.Organization + "/" + name + ".git"

    if r, found := g.github_source.Repos[name] ; ! found {
        r = RepositoryStruct{}
		r.set(repo, map[string]string {"origin":upstream}, map[string]string {"master":"origin/master"})
        g.github_source.Repos[name] = r
        return true // New added
    }
    return false
}

func (g *GitHubStruct)save_yaml(file string) error {

    d, err := yaml.Marshal(&g.github_source)
    if  err != nil {
        return fmt.Errorf("Unable to encode github data in yaml. %s", err)
    }

    if err := ioutil.WriteFile(file, d, 0644) ; err != nil {
        return fmt.Errorf("Unable to save '%s'. %s", file, err)
    }
    return nil
}

func (g *GitHubStruct)load_yaml(file string) error {
    d, err := ioutil.ReadFile(file)
    if err != nil {
        return fmt.Errorf("Unable to load '%s'. %s", file, err)
    }

    err = yaml.Unmarshal(d, &g.github_source)
    if  err != nil {
        return fmt.Errorf("Unable to decode github data in yaml. %s", err)
    }
    return nil
}

func (r *CreateArgReq) SaveMaintainOptions(ret *goforjj.PluginData) {
	if ret.Options == nil {
		ret.Options = make(map[string]goforjj.PluginOption)
	}
}
