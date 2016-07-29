package main

//go:generate go get github.hpe.com/christophe-larsonneur/goforjj gopkg.in/yaml.v2
//go:generate go build -o $GOPATH/bin/forjj-genapp github.hpe.com/christophe-larsonneur/goforjj/genapp
//go:generate forjj-genapp github.yaml

import (
    "github.hpe.com/christophe-larsonneur/goforjj"
    "github.com/google/go-github/github"
)

type GitHubStruct struct {
    source string
    workspace string
    token string
    debug bool
    Client *github.Client
}

type GitHubSourceStruct struct {
    Urls goforjj.PluginService           // github base Url
    Repos map[string]RepositoryStruct    // Collection of repositories managed in github
    Organization string                  // Organization name
    UserGroups []UserGroupStruct         // Collection of groups to add to the organization
}

type UserGroupStruct struct {
    Name string // Name of the group
    Role string // Role to apply in the context
}

type RepositoryStruct  struct {
    goforjj.PluginRepo           // Name/Upstream
    Description string           // Title in github repository
    UserGroups []UserGroupStruct // Collection of groups to add to the organization
}

const github_source_file = "github.yaml"
