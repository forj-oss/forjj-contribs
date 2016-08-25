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
    UserGroups []UserGroupStruct         // Collection of groups to add to the organization
}

type UserGroupStruct struct {
    Name string // Name of the group
    Role string // Role to apply in the context
}

type RepositoryStruct  struct {
    goforjj.PluginRepo `,inline`      // Name/Upstream
    Description string                // Title in github repository
    UserGroups []UserGroupStruct      // Collection of groups to add to the organization
}

const github_source_file = "github.yaml"
