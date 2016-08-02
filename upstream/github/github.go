package main

import (
    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
    "github.hpe.com/christophe-larsonneur/goforjj"
    "fmt"
    "log"
    "net/url"
    "regexp"
)

func (g *GitHubStruct)github_connect(server string, ret *goforjj.PluginData) (* github.Client) {
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: g.token})
    tc := oauth2.NewClient(oauth2.NoContext, ts)

    g.Client = github.NewClient(tc)

    if err := g.github_set_url(server) ; err != nil {
        ret.ErrorMessage = fmt.Sprintf("Invalid url. %s", err)
        return nil
    }
    log.Printf("Github Base URL used : %s", g.Client.BaseURL)

    if user , _, err := g.Client.Users.Get("") ; err != nil {
        ret.ErrorMessage = fmt.Sprintf("Unable to get the owner of the token given. %s", err)
        return nil
    } else {
        log.Printf("Connection successful. Token given to user '%s'", *user.Login)
    }

    return  g.Client
}

func (g *GitHubStruct)github_set_url(server string) (err error) {
    if server == ""  || server == "github.com" || server == "https://github.com"{
        return
    }

    if found, _ := regexp.MatchString("^https?://.*", server) ; found {
        g.Client.BaseURL, err = url.Parse(server)
        if err != nil {
            return
        }
    } else {
        g.Client.BaseURL.Host = server
    }

    if g.Client.BaseURL.Path == "/" {
        g.Client.BaseURL.Path = "/api/v3/"
    }

    if g.Client.BaseURL.Scheme == "" {
        g.Client.BaseURL.Scheme = "https"
    }
    return
}
