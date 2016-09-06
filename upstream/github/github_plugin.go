package main

import (
    "github.hpe.com/christophe-larsonneur/goforjj"
    "github.com/google/go-github/github"
)

type GitHubStruct struct {
    source_mount string              // mount point
    workspace_mount string           // mount point
    token string
    debug bool
    user string                      // github user name
    Client *github.Client
    github_source GitHubSourceStruct // github source structure (yaml)
}

type GitHubSourceStruct struct {
    goforjj.PluginService `,inline`      // github base Url
    Repos map[string]RepositoryStruct    // Collection of repositories managed in github
    Organization string                  // Organization name
    OrgDisplayName string                // Organization's display name.
    Users map[string]string              // Collection of users role at organization level
    Groups map[string]string             // Collection of groups role at organization level
}

type RepositoryStruct  struct { // Used to stored the yaml source file. Not used to respond to the API requester.
    Name string                 // Name of the Repo
    Flow string                 // Flow applied on the repo.
    Description string          // Title in github repository
    Users map[string]string     // Collection of users role
    Groups map[string]string     // Collection of groups role
    // Following data are used at runtime but not saved. Used to respond to the API.
    exist bool                      // True if the repo exist.
    remotes map[string]string       // k: remote name, v: remote url
    branchConnect map[string]string // k: local branch name, v: remote/branch
}

const github_source_file = "github.yaml"
