package main

import (
    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
    "github.hpe.com/christophe-larsonneur/goforjj"
    "fmt"
    "log"
)

func (g *GitHubStruct)github_connect(ret *goforjj.PluginData) (* github.Client) {
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: g.token})
    tc := oauth2.NewClient(oauth2.NoContext, ts)

    g.Client = github.NewClient(tc)

    if user , _, err := g.Client.Users.Get("") ; err != nil {
        ret.ErrorMessage = fmt.Sprintf("Unable to get the owner of the token given. %s", err)
        return nil
    } else {
        log.Printf("Connection successful. Token given to user '%s'", user.Login)
    }

    return  g.Client
}
