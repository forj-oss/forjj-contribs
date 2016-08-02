package main

import (
    "gopkg.in/yaml.v2"
    "fmt"
    "io/ioutil"
)

func (g *GitHubStruct)create_yaml(file string, req *CreateReq) error {
    // Write the github.yaml source file.
    var g.source GitHubSourceStruct
    g.source.Url = make(map[string]string)
    g.source.Url["github-base-url"] = g.Client.BaseURL.String()

    if orga := req.GithubOrganization; orga == "" {
        g.source.Organization = req.ForjjOrganization
    } else {
        g.source.Organization = req.GithubOrganization
    }

    // Ensure Infra is already in the list of repo managed.
    if g.source.Repos == nil {
        g.source.Repos = make(map[string]RepositoryStruct)
    }

    infra, found := g.source.Repos[req.ForjjInfra]
    if ! found {
        infra = RepositoryStruct{
            Name: req.ForjjInfra,
            Upstream: "git@" + g.Client.BaseURL.Host + ":" + g.source.orga + "/" + req.ForjjInfra + ".git",
            Description: fmt.Sprintf("Infrastructure repository for Organization '%s' maintained by Forjj", g.source.orga),
            UserGroups: make([]UserGroupStruct, 0),
        }
    }


    d, err := yaml.Marshal(&g.source)
    if  err != nil {
        return fmt.Errorf("Unable to encode github data in yaml. %s", err)
    }

    if err := ioutil.WriteFile(file, d, 0644) ; err != nil {
        return fmt.Errorf("Unable to save 'github.yaml'. %s", err)
    }
    return nil
}
