package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"regexp"

	"github.com/forj-oss/goforjj"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"strings"
	"strconv"
)

func (req *CreateReq) InitOrganization(g *GitHubStruct) {
	instance := req.Forj.ForjjInstanceName
	if orga := req.Objects.App[instance].Organization; orga == "" {
		g.github_source.Organization = req.Objects.App[instance].ForjjOrganization
	} else {
		g.github_source.Organization = orga
	}

}

// No change for now.
func (req *UpdateReq) InitOrganization(g *GitHubStruct) {
}

func (g *GitHubStruct) github_connect(server string, ret *goforjj.PluginData) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: g.token})
	g.ctxt = context.Background()
	tc := oauth2.NewClient(g.ctxt, ts)

	g.Client = github.NewClient(tc)

	if err := g.github_set_url(server); err != nil {
		ret.Errorf("Invalid url. %s", err)
		return nil
	}

	log.Printf("Github API URL used : %s", g.Client.BaseURL)
	log.Printf("Github URL used : %s", g.github_source.Urls["github-url"])

	if user, _, err := g.Client.Users.Get(g.ctxt, ""); err != nil {
		ret.Errorf("Unable to get the owner of the token given. %s", err)
		return nil
	} else {
		g.user = *user.Login
		log.Printf("%s. Token given by user '%s'", ret.StatusAdd("Connection successful."), *user.Login)
	}

	return g.Client
}

// github_set_url will set Urls/github-base-url in create/update context
// and set the url object when base-url is not empty (private GitHub)
func (g *GitHubStruct) github_set_url(server string) (err error) {
	gh_url := ""
	if g.github_source.Urls == nil {
		g.github_source.Urls = make(map[string]string)
	}
	if !g.maintain_ctxt {
		if server == "" || server == "api.github.com" || server == "github.com" {
			g.github_source.Urls["github-base-url"] = "https://api.github.com/" // Default public API link
			g.github_source.Urls["github-url"] = "https://github.com"          // Default public link
		} else {
			// To accept GitHub entreprise without ssl, permit server to have url format.
			var entr_github_re *regexp.Regexp
			if re, err := regexp.Compile("^(https?://)(.*)(/api/v3)/?$"); err != nil {
				return err
			} else {
				entr_github_re = re
			}
			res := entr_github_re.FindAllString(server, -1)
			if res == nil {
				gh_url = "https://" + server + "/api/v3/"
				g.github_source.Urls["github-url"] = "https://" + server
			} else {
				if res[2] == "" {
					return fmt.Errorf("Unable to determine github URL from '%s'. It must be [https?://]Server[:Port][/api/v3]", server)
				}
				if res[1] == "" {
					gh_url += "https://"
				}
				gh_url += res[2]
				g.github_source.Urls["github-url"] = gh_url
				gh_url += "/api/v3/"
			}
			g.github_source.Urls["github-base-url"] = gh_url
		}
	} else {
		gh_url = g.github_source.Urls["github-base-url"]
	}

	if gh_url == "" {
		return
	}

	g.Client.BaseURL, err = url.Parse(gh_url)
	if err != nil {
		return
	}

	/*	// Adding api/V3 for server given or url without path, ie http?://<server> instead or http?://<server>/<path>?
		if g.Client.BaseURL.Path == "" {
			log.Printf("Adding /api/v3 to github url given %s", gh_url)
			g.Client.BaseURL.Path = "/api/v3/"
			g.github_source.Urls["github-base-url"] = g.Client.BaseURL.String()
		}*/
	return
}

type GithubEntrepriseOrganization struct {
	Login        string
	Profile_name string
	Admin        string
}

// Ensure organization exists means:
// - organization exist. if not it is created.
// - organization has current user as owner
func (g *GitHubStruct) ensure_organization_exists(ret *goforjj.PluginData) (s bool) {

	s = false
	_, resp, err := g.Client.Organizations.Get(g.ctxt, g.github_source.Organization)
	if err != nil && resp == nil {
		log.Printf(ret.Errorf("Unable to get '%s' organization information. %s", g.github_source.Organization, err))
		return
	}
	if resp.StatusCode != 200 {
		// need to create the Organization
		var orga GithubEntrepriseOrganization = GithubEntrepriseOrganization{g.github_source.Organization, g.github_source.OrgDisplayName, g.user}
		var res_orga github.Organization

		if v, found := g.github_source.Urls["github-base-url"]; !found || v == "" {
			log.Printf(ret.StatusAdd("Unable to create an organization on github.com. You must do it manually."))
			return false
		}

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

// setOrganizationTeams maintain the list of teams as defined by the github.yaml file.
func (g *GitHubStruct) setOrganizationTeams(ret *goforjj.PluginData) (_ bool) {
	if g.github_source.NoTeams {
		log.Printf(ret.StatusAdd("Users/Groups maintain ignored."))
		return true
	}
	// Load teams list already defined in github.
	github_teams, resp, err := g.Client.Organizations.ListTeams(g.ctxt, g.github_source.Organization, nil)
	if err != nil && resp == nil {
		log.Printf(ret.Errorf("Unable to verify '%s' organization teams. %s", g.github_source.Organization, err))
		return
	}

	teams := make(map[string]*github.Team)
	// Loop on github list to find out undefined team in github.yaml, then remove them
	if resp.StatusCode == 200 {
		for _, github_team := range github_teams {
			// TODO: Support more teams options to maintain
			if _, found := g.github_source.Groups[*github_team.Name]; !found {
				// Remove team
				if g.new_forge && ! g.force {
					ret.Errorf("Unable to remove teams on a new Forge if you do not forcelly request it. " +
						"To fix it, use the github force option or update your Forjfile.")
					return
				}
				log.Printf(ret.StatusAdd("Removing uncontrolled team '%s'.", *github_team.Name))
				resp, err = g.Client.Organizations.DeleteTeam(g.ctxt, *github_team.ID)
				if err != nil && resp == nil {
					log.Printf(ret.Errorf("Unable to remove team '%s' from '%s' organization. %s",
						*github_team.Name, g.github_source.Organization, err))
					return
				} else if resp.StatusCode != 204 {
					log.Printf(ret.Errorf("Unable to remove team '%s' from '%s' organization. HTTP status : %s",
						*github_team.Name, g.github_source.Organization, resp.Status))
					return

				}
			} else {
				teams[*github_team.Name] = github_team
			}
		}
	}

	valid_roles := map[string]string{"admin": "admin", "push": "push", "pull": "pull"}

	// Loop on github.yaml group to create/update teams
	for name, team := range g.github_source.Groups {
		if github_team, found := teams[name]; found {
			updated := false
			team_to_update := github.Team{Name: github_team.Name}
			if _, valid := valid_roles[team.Role]; !valid {
				team.Role = "pull"
			}
			if team.Role != *github_team.Permission {
				team_to_update.Permission = &team.Role
				updated = true
			}
			if updated {
				log.Printf(ret.StatusAdd("Updating team '%s'.", name))

				_, resp, err = g.Client.Organizations.EditTeam(g.ctxt, *github_team.ID, &team_to_update)
				if err != nil && resp == nil {
					log.Printf(ret.Errorf("Unable to update organization team '%s'. %s", name, err))
					return
				}
				if resp.StatusCode != 200 {
					log.Printf(ret.Errorf("Unable to update organization team '%s'. %s", name, resp.Status))
					return
				}
			} else {
				log.Printf(ret.StatusAdd("No change on team '%s'.", name))
			}
			g.setOrganizationTeamsMembers(ret, github_team)
			continue
		}

		// Team have to be created
		log.Printf(ret.StatusAdd("Creating team '%s'.", name))
		github_team := new(github.Team)
		github_team.Name = &name
		if _, valid := valid_roles[team.Role]; valid {
			github_team.Permission = &team.Role
		} else {
			github_team.Permission = nil
		}
		github_team, resp, err = g.Client.Organizations.CreateTeam(g.ctxt, g.github_source.Organization, github_team)
		if err != nil && resp == nil {
			log.Printf(ret.Errorf("Unable to create organization team '%s'. %s", name, err))
			return
		}
		if resp.StatusCode != 201 {
			log.Printf(ret.Errorf("Unable to create organization team '%s'. %s", name, resp.Status))
			return
		}

		g.setOrganizationTeamsMembers(ret, github_team)
	}
	return true
}

func (g *GitHubStruct) setOrganizationTeamsMembers(ret *goforjj.PluginData, team *github.Team) (_ bool) {
	github_users, resp, err := g.Client.Organizations.ListTeamMembers(g.ctxt, *team.ID, nil)
	if err != nil && resp == nil {
		log.Printf(ret.Errorf("Unable to check team '%s'. %s", team.Name, err))
		return false
	}
	users := make(map[string]int)
	var team_source TeamStruct
	if t, found := g.github_source.Groups[*team.Name]; !found {
		log.Printf(ret.StatusAdd("Warning. team '%s' has no membership declared", *team.Name))
		return false
	} else {
		team_source = t
	}

	// GetTeamMembers and remove those missing in github.yaml
	for _, user := range github_users {
		if found, _ := goforjj.InArray(*user.Login, team_source.Users); !found {
			log.Printf(ret.StatusAdd("Removing unknown user '%s'.", *user.Login))
			resp, err = g.Client.Organizations.RemoveTeamMembership(g.ctxt, *team.ID, *user.Login)
			if err != nil && resp == nil {
				log.Printf(ret.Errorf("Unable to remove team member '%s' from team '%s'. %s",
					*user.Login, *team.Name, err))
				return
			}
			if resp.StatusCode != 204 {
				log.Printf(ret.Errorf("Unable to remove team member '%s' from team '%s'. %s",
					*user.Name, *team.Name, resp.Status))
				return
			}
			continue
		}
		users[*user.Login] = 1 // Use map key facility only
	}

	// TODO: Detect new members to add
	for _, user := range team_source.Users {
		if _, found := users[user]; !found {
			log.Printf(ret.StatusAdd("Adding missing user '%s'.", user))
			_, resp, err = g.Client.Organizations.AddTeamMembership(g.ctxt, *team.ID, user, nil)
			if err != nil && resp == nil {
				log.Printf(ret.Errorf("Unable to add team member '%s' to team '%s'. %s",
					user, *team.Name, err))
				return
			}
			if resp.StatusCode != 200 {
				log.Printf(ret.Errorf("Unable to add team member '%s' to team '%s'. %s",
					user, *team.Name, resp.Status))
				return
			}
		}
	}
	return true
}

// Return an error if at least one repo exist. Used at create/update time.
func (g *GitHubStruct) repos_exists(ret *goforjj.PluginData) (err error) {
	c := g.Client.Repositories

	// loop on list of repos, and ensure they exist with minimal config and rights
	for name, repo_data := range g.github_source.Repos {
		if found_repo, _, e := c.Get(g.ctxt, g.github_source.Organization, name); e == nil {
			if err == nil && name == g.app.ForjjInfra { // Infra repository.
				err = fmt.Errorf("Infra repository '%s' already exist in github server.", name)
			}
			repo_data.exist = true
			if repo_data.remotes == nil {
				repo_data.remotes = make(map[string]goforjj.PluginRepoRemoteUrl)
				repo_data.branchConnect = make(map[string]string)
			}
			repo_data.remotes["origin"] = goforjj.PluginRepoRemoteUrl{
				Ssh: *found_repo.SSHURL,
				Url: *found_repo.HTMLURL,
			}
			repo_data.branchConnect["master"] = "origin/master"
		}
		ret.Repos[name] = goforjj.PluginRepo{
			Name:          repo_data.Name,
			Exist:         repo_data.exist,
			Remotes:       repo_data.remotes,
			BranchConnect: repo_data.branchConnect,
			Owner:         g.github_source.Organization,
		}
	}
	return
}

func (g *GitHubStruct) IsNewForge(ret *goforjj.PluginData) (_ bool) {
	c := g.Client.Repositories

	// loop on list of repos, and ensure they exist with minimal config and rights
	for name, repo := range g.github_source.Repos {
		if ! repo.Infra {
			continue
		}
		// Infra repository.
		if _, resp, e := c.Get(g.ctxt, g.github_source.Organization, name); e != nil && resp == nil {
			ret.Errorf("Unable to identify the infra repository. Unknown issue: %s", e)
			return
		} else {
			g.new_forge = (resp.StatusCode != 200)
		}
		return true
	}
	ret.Errorf("Unable to identify the infra repository. At least, one repo must be identified with " +
		"`%s` in %s. You can use Forjj update to fix this.", "Infra: true", "github")
	return
}


// Populate ret.Repos with req.repos status and information
func (g *GitHubStruct) req_repos_exists(req *UpdateReq, ret *goforjj.PluginData) (err error) {
	if req == nil || ret == nil {
		return fmt.Errorf("Internal error: Invalid parameters. req and ret cannot be nil.")
	}

	c := g.Client.Repositories

	// loop on list of repos, and ensure they exist with minimal config and rights
	for name, _ := range req.Objects.Repo {
		log.Printf("Looking for Repo '%s' from '%s'", name, g.github_source.Organization)
		found_repo, _, err := c.Get(g.ctxt, g.github_source.Organization, name)

		r := goforjj.PluginRepo{
			Name:          name,
			Exist:         (err == nil),
			Remotes:       make(map[string]goforjj.PluginRepoRemoteUrl),
			BranchConnect: make(map[string]string),
			Owner:         *found_repo.Organization.Name,
		}
		if err == nil {
			r.Remotes["origin"] = goforjj.PluginRepoRemoteUrl{
				Ssh: *found_repo.SSHURL,
				Url: *found_repo.HTMLURL,
			}
			r.BranchConnect["master"] = "origin/master"
		}

		ret.Repos[name] = r
	}
	return
}

func (r *RepositoryStruct) exists(gws *GitHubStruct) bool {
	c := gws.Client.Repositories
	_, _, err := c.Get(gws.ctxt, gws.github_source.Organization, r.Name)

	if err == nil { // repos exist
		return true
	}
	return false
}

// FUTURE: Add users/groups

func (r *RepositoryStruct) ensure_exists(gws *GitHubStruct, ret *goforjj.PluginData) error {
	// test existence
	c := gws.Client.Repositories
	found_repo, _, err := c.Get(gws.ctxt, gws.github_source.Organization, r.Name)
	if err != nil {
		// Creating repository
		github_repo := github.Repository{
			Description: &r.Description,
			Name:        &r.Name,
			HasIssues:   &r.IssueTracker,
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
	repo, found := ret.Repos[r.Name]
	if !found {
		repo = goforjj.PluginRepo{
			Name:          r.Name,
			Remotes:       make(map[string]goforjj.PluginRepoRemoteUrl),
			Exist:         true,
			BranchConnect: make(map[string]string),
		}
	}

	// TODO: See how to integrate the flow change here to respond the proper branch connect.
	repo.BranchConnect["master"] = "origin/master"
	if found_repo.Parent != nil {
		repo.Remotes["upstream"] = goforjj.PluginRepoRemoteUrl{
			Ssh: *found_repo.Parent.SSHURL,
			Url: *found_repo.Parent.HTMLURL,
		}
	}

	repo.Remotes["origin"] = goforjj.PluginRepoRemoteUrl{
		Ssh: *found_repo.SSHURL,
		Url: *found_repo.HTMLURL,
	}

	repo.Owner = gws.github_source.Organization

	ret.Repos[r.Name] = repo
	return nil
}

func (r *RepositoryStruct) maintain(e_repo *github.Repository) *github.Repository {
	if e_repo == nil {
		return nil
	}
	update := false
	ret := github.Repository{}
	ret.Name = e_repo.Name
	update = update || updateString(e_repo.Description, &ret.Description, r.Description, "Description")
	update = update || updateBool(e_repo.HasIssues, &ret.HasIssues, r.IssueTracker, "Issue tracker")

	if update {
		return &ret
	}
	return nil

}

func updateString(orig *string, dest **string, to, field string) (updated bool) {
	var from string

	defer func() {
		if updated {
			log.Printf("%s: %s => %s", field, from, to)
			*dest = &to
		}
	}()

	if orig != nil {
		from = *orig
	}
	if from != to {
		updated = true
	}
	return
}

func updateBool(orig *bool, dest **bool, to bool, field string) (updated bool) {
	var from bool
	defer func() {
		if updated {
			log.Printf("%s: %t => %t", field, from, to)
			*dest = &to
		}
	}()

	if orig != nil {
		from = *orig
	}
	if from != to {
		updated = true
	}
	return
}

func (g *GitHubStruct) SetOrgHooks(org_hook_disabled, repo_hook_disabled, wh_policy string, hooks map[string]WebhooksInstanceStruct) {

	if b, err := strconv.ParseBool(org_hook_disabled) ; err != nil {
		log.Printf("Organization webhook disabled: invalid boolean: %s", org_hook_disabled)
		g.github_source.NoOrgHook = b
		g.github_source.WebHooks = make(map[string]WebHookStruct)
	}

	if b, err := strconv.ParseBool(repo_hook_disabled); err != nil {
		log.Printf("Organization webhook disabled: invalid boolean: %s", repo_hook_disabled)
		g.github_source.NoRepoHook = b
	}

	if v := inStringList(wh_policy, "manage", "sync"); v == "" {
		if wh_policy != "" {
			log.Printf("'Invalid value '%s' for 'WebhooksManagement'. Set it to 'sync'.", wh_policy)
		} else {
			log.Print("'WebhooksManagement' is set by default to 'sync'.")
		}
		g.github_source.WebHookPolicy = "sync"
	} else {
		g.github_source.WebHookPolicy = v
	}

	if g.github_source.NoOrgHook {
		return
	}

	for name, hook := range hooks {
		if hook.Organization == "false" {
			continue
		}

		data := WebHookStruct{
			Url: hook.Url,
			Events: strings.Split(hook.Events, ","),
			Enabled: hook.Enabled,
		}
		if v, err := strconv.ParseBool(hook.SslCheck); err == nil {
			data.SSLCheck = v
			log.Printf("SSL Check '%s' => %t", name, v)
		} else {
			log.Printf("SSLCheck has an invalid boolean string representation '%s'. Ignored. SSL Check is set to true.",
				name)
			data.SSLCheck = true
		}

		g.github_source.WebHooks[name] = data
	}
}
