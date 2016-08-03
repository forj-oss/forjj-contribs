package main

import (
    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
    "github.hpe.com/christophe-larsonneur/goforjj"
    "log"
    "net/url"
    "regexp"
)

func (g *GitHubStruct)github_connect(server string, ret *goforjj.PluginData) (* github.Client) {
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: g.token})
    tc := oauth2.NewClient(oauth2.NoContext, ts)

    g.Client = github.NewClient(tc)

    if err := g.github_set_url(server) ; err != nil {
        ret.Errorf("Invalid url. %s", err)
        return nil
    }
    log.Printf("Github Base URL used : %s", g.Client.BaseURL)

    if user , _, err := g.Client.Users.Get("") ; err != nil {
        ret.Errorf("Unable to get the owner of the token given. %s", err)
        return nil
    } else {
        g.user = *user.Login
        log.Printf("%s. Token given by user '%s'", ret.StatusAdd("Connection successful."), *user.Login)
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

type GithubEntrepriseOrganization struct {
    Login string
    Profile_name string
    Admin string
}

// Ensure organization exists means:
// - organization exist. if not it is created.
// - organization has current user as owner
func (g *GitHubStruct)ensure_organization_exists(ret *goforjj.PluginData) (s bool) {

    s = false
    _, resp, err := g.Client.Organizations.Get(g.github_source.Organization)
    if err != nil && resp == nil {
        log.Printf(ret.Errorf("Unable to get '%s' organization information. %s", g.github_source.Organization, err))
        return
    }
    if resp.StatusCode != 200 {
        // need to create the Organization
        var orga GithubEntrepriseOrganization = GithubEntrepriseOrganization{ g.github_source.Organization, g.github_source.OrgDisplayName, g.user }
        var res_orga github.Organization

        req, err := g.Client.NewRequest("POST", "admin/organizations", orga)
        if err != nil {
            log.Printf(ret.Errorf("Unable to create '%s' as organization. Request is failing. %s", g.github_source.Organization, err))
            return
        }

        _, err = g.Client.Do(req, res_orga)
        if err != nil {
            log.Printf(ret.Errorf("Unable to create '%s' as organization. %s.\nYour credentials is probably insufficient.\nYou can update your token access rights or ask to create the organization and attach a Full control access token to the organization owner dedicated to Forjj.\nAs soon as fixed, your can restart forjj maintain", g.github_source.Organization, err))
            return
        }
        _, resp, err = g.Client.Organizations.Get(g.github_source.Organization)
        if err != nil && resp == nil {
            log.Printf(ret.Errorf("Unable to get '%s' organization information. %s", g.github_source.Organization, err))
            return
        }
        if resp.StatusCode != 200 {
            log.Printf(ret.Errorf("Unable to get '%s' created organization information. %s", g.github_source.Organization, err))
            return
        }
        log.Printf(ret.StatusAdd("'%s' organization created", g.github_source.Organization))
    } else {
        // Ensure the organization is writable
        _, resp, err := g.Client.Organizations.IsMember(g.github_source.Organization, g.user)
        if err != nil && resp == nil {
            log.Printf(ret.Errorf("Unable to verify '%s' organization ownership. %s", g.github_source.Organization, err))
            return
        }
        if resp.StatusCode == 302 {
            log.Printf(ret.Errorf("'%s' organization is not owned by '%s'. This is a Forjj requirement. Please ask the owner to add '%s' as owner of this organization.", g.github_source.Organization, g.user, g.user))
            return
        }
        log.Printf(ret.StatusAdd("'%s' organization access verified", g.github_source.Organization))
    }
    return true
}
