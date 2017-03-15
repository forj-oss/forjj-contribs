package main

import (
    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
    "github.com/forj-oss/goforjj"
    "log"
    "net/url"
    "regexp"
    "fmt"
    "context"
)

func (req *CreateReq)InitOrganization(g *GitHubStruct) {
    instance := req.Forj.ForjjInstanceName
    if orga := req.Objects.App[instance].Add.Organization; orga == "" {
        g.github_source.Organization = req.Objects.App[instance].Add.ForjjOrganization
    } else {
        g.github_source.Organization = orga
    }

}

// No change for now.
func (req *UpdateReq)InitOrganization(g *GitHubStruct) {
}

func (g *GitHubStruct)github_connect(server string, ret *goforjj.PluginData) (* github.Client) {
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: g.token})
    g.ctxt = context.Background()
    tc := oauth2.NewClient(g.ctxt, ts)

    g.Client = github.NewClient(tc)


    if du, found := g.github_source.Urls["github-base-url"] ; !found || (found && du == "") || server != "" {
        if err := g.github_set_url(server) ; err != nil {
            ret.Errorf("Invalid url. %s", err)
            return nil
        }
    } else {
        log.Printf("Using github-base-url : %s", du)
        if u, err := url.Parse(du) ; err != nil {
            return nil
        } else {
            g.Client.BaseURL = u
        }
    }
    log.Printf("Github Base URL used : %s", g.Client.BaseURL)

    if user , _, err := g.Client.Users.Get(g.ctxt, "") ; err != nil {
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
    _, resp, err := g.Client.Organizations.Get(g.ctxt, g.github_source.Organization)
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

        _, err = g.Client.Do(g.ctxt, req, res_orga)
        if err != nil {
            log.Printf(ret.Errorf("Unable to create '%s' as organization. %s.\nYour credentials is probably insufficient.\nYou can update your token access rights or ask to create the organization and attach a Full control access token to the organization owner dedicated to Forjj.\nAs soon as fixed, your can restart forjj maintain", g.github_source.Organization, err))
            return
        }
        _, resp, err = g.Client.Organizations.Get(g.ctxt, g.github_source.Organization)
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
        _, resp, err := g.Client.Organizations.IsMember(g.ctxt, g.github_source.Organization, g.user)
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

// Return an error if at least one repo exist. Used at create/update time.
func (g *GitHubStruct)repos_exists(ret *goforjj.PluginData) (err error) {
    c := g.Client.Repositories

    // loop on list of repos, and ensure they exist with minimal config and rights
    for name, repo_data := range g.github_source.Repos {
        if found_repo, _, e := c.Get(g.ctxt, g.github_source.Organization, name) ; e == nil {
            if err == nil {
                err = fmt.Errorf("At least '%s' already exist in github server.", name)
            }
            repo_data.exist = true
            if repo_data.remotes == nil {
                repo_data.remotes = make(map[string]string)
                repo_data.branchConnect = make(map[string]string)
            }
            repo_data.remotes["origin"] = *found_repo.SSHURL
            repo_data.branchConnect["master"] = "origin/master"
        }
        if ret != nil {
            ret.Repos[name] = goforjj.PluginRepo{
                Name: repo_data.Name,
                Exist: repo_data.exist,
                Remotes: repo_data.remotes,
                BranchConnect: repo_data.branchConnect,
            }
        }
    }
    return
}

// Populate ret.Repos with req.repos status and information
func (g *GitHubStruct)req_repos_exists(req *UpdateReq, ret *goforjj.PluginData) (err error) {
    if req == nil || ret == nil {
        return fmt.Errorf("Internal error: Invalid parameters. req and ret cannot be nil.")
    }

    c := g.Client.Repositories

    // loop on list of repos, and ensure they exist with minimal config and rights
    for name, _ := range req.Objects.Repo {
        log.Printf("Looking for Repo '%s' from '%s'", name, g.github_source.Organization)
        found_repo, _, err := c.Get(g.ctxt, g.github_source.Organization, name)

        r := goforjj.PluginRepo{
            Name: name,
            Exist: (err == nil),
            Remotes: make(map[string]string),
            BranchConnect: make(map[string]string),
        }
        if err == nil {
            r.Remotes["origin"] = *found_repo.SSHURL
            r.BranchConnect["master"] = "origin/master"
        }

        ret.Repos[name] = r
    }
    return
}

func (r *RepositoryStruct)exists(gws *GitHubStruct) bool{
    c := gws.Client.Repositories
    _, _, err := c.Get(gws.ctxt, gws.github_source.Organization, r.Name)

    if err == nil { // repos exist
        return true
    }
    return false
}

// FUTURE: Add users/groups

func (r *RepositoryStruct)ensure_exists(gws *GitHubStruct, ret *goforjj.PluginData) error{
     // test existence
    c := gws.Client.Repositories
    found_repo, _, err := c.Get(gws.ctxt, gws.github_source.Organization, r.Name)
    if err != nil {
        // Creating repository
        github_repo := github.Repository{
            Description: &r.Description,
            Name: &r.Name,
        }
        found_repo, _, err = c.Create(gws.ctxt, gws.github_source.Organization, &github_repo)
        if err != nil {
            ret.Errorf("Unable to create '%s' in organization '%s'. %s.", r.Name, gws.github_source.Organization, err)
            return err
        }
        log.Printf(ret.StatusAdd("Repo '%s': created", r.Name))

    } else {
        // Updating repository if needed
        repo_updated := r.maintain(found_repo)
        if repo_updated == nil {
            log.Printf(ret.StatusAdd("Repo '%s': No change", r.Name))
        } else {
            found_repo, _, err = c.Edit(gws.ctxt, gws.github_source.Organization, r.Name, repo_updated)
            if err != nil {
                ret.Errorf("Unable to update '%s' in organization '%s'. %s.", r.Name, gws.github_source.Organization, err)
                return err
            }
            log.Printf(ret.StatusAdd("Repo '%s': updated", r.Name))
        }
    }

    // TODO: Use a goforjj function to manage this return.

    // Prepare return status information to github API caller.
    if ret.Repos == nil {
        ret.Repos = make(map[string]goforjj.PluginRepo)
    }

    // TODO: Add github flow driver for repos management
    if repo, found := ret.Repos[r.Name]; found {
        repo.Remotes["origin"] = *found_repo.SSHURL
        ret.Repos[r.Name] = repo
    } else {
        repo = goforjj.PluginRepo {
            Name: r.Name,
            Remotes: make(map[string]string),
            Exist: true,
            BranchConnect: make(map[string]string),
        }

        // TODO: See how to integrate the flow change here to respond the proper branch connect.
        repo.Remotes["origin"] = *found_repo.SSHURL
        repo.BranchConnect["master"] = "origin/master"
        if found_repo.Parent != nil {
            repo.Remotes["upstream"] = *found_repo.Parent.HTMLURL
        }
        ret.Repos[r.Name] = repo
    }
    return nil
}

func (r *RepositoryStruct)maintain(e_repo *github.Repository) *github.Repository {
    update := false
    ret := github.Repository{}
    ret.Name = e_repo.Name
    if e_repo.Description != &r.Description {
        update = true
        ret.Description = &r.Description
    }

    if update {
        return &ret
    }
    return nil

}
