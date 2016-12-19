package main

import (
    "github.com/forj-oss/goforjj"
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

const github_source_file = "github.yaml"
