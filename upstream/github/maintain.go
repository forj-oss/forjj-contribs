// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"github.com/forj-oss/goforjj"
	"log"
	"github.com/google/go-github/github"
)

func (g *GitHubStruct)MaintainOrgHooks(ret *goforjj.PluginData) (_ bool) {
	// organization level
	if hooks, _, err := g.Client.Organizations.ListHooks(g.ctxt, g.github_source.Organization, nil) ; err == nil {
		for _, hook := range hooks {
			if hook.GetName() != "web" {
				continue
			}
			if h, found := g.GetWebHook(hook, g.github_source.WebHooks) ; found {
				h.identified = true
				g.github_source.WebHooks[h.name] = h
				if h.Update(hook) {
					if _, _, err := g.Client.Organizations.EditHook(g.ctxt, g.github_source.Organization, hook.GetID(), hook) ; err != nil {
						log.Print(ret.Errorf("Failed to update '%s'. %s", hook.GetName(), err))
						return
					} else {
						log.Print(ret.StatusAdd("WebHook '%s' updated.", h.name))
					}
				}
			} else if g.github_source.WebHookPolicy == "" || g.github_source.WebHookPolicy == "sync" {
				if _, err := g.Client.Organizations.DeleteHook(g.ctxt, g.github_source.Organization, hook.GetID()); err != nil {
					log.Print(ret.Errorf("Failed to delete '%s'. %s", hook.Config["url"], err))
					return
				} else {
					log.Print(ret.StatusAdd("Org WebHook '%s' removed.", hook.Config["url"]))
				}
			} else {
				log.Print(ret.StatusAdd("Org WebHook '%s' not managed. (org webhook policy = 'manage')", hook.Config["url"]))
			}
		}
	}
	for name, hook := range g.github_source.WebHooks {
		if hook.identified {
			continue
		}

		Config := make(map[string]interface{})

		Config["url"] = hook.Url
		Config["insecure_ssl"] = g.SetWebHookInsecure(hook.SSLCheck)
		Config["content_type"] = hook.ContentType

		integ_name := "web" // webhook integration type in a dedicated webhook tab.
		new_hook := github.Hook{
			Name: &integ_name,
			Config: Config,
			Events: hook.Events,
		}

		hook.HookEnabled(&new_hook)

		if _, resp, err := g.Client.Organizations.CreateHook(g.ctxt, g.github_source.Organization, &new_hook) ; err != nil {
			log.Print(ret.Errorf("Failed to create '%s'. %s", name, err))
		} else if resp.StatusCode != 201 {
			log.Print(ret.Errorf("Failed to create '%s'. %s", name, resp.Status))
			return
		}
		log.Print(ret.StatusAdd("Org WebHook '%s created.", name))
		hook.identified = true
	}
	return true
}

func (GitHubStruct)SetWebHookInsecure(insecure bool) string {
	if ! insecure {
		return "1"
	}
	return "0"
}

func (GitHubStruct)GetWebHook(hook *github.Hook, webhooks map[string]WebHookStruct) (_ WebHookStruct, _ bool) {
	hook_url := ""
	if url, found := hook.Config["url"]; found {
		if u, ok := url.(string); ok {
			hook_url = u
		}
	}
	for name, h := range webhooks {
		if h.Url != hook_url {
			continue
		}
		h.name = name
		return h, true
	}
	return
}

func (g *GitHubStruct) MaintainHooks(repo *RepositoryStruct, ret *goforjj.PluginData) (_ bool) {
	// Repository level
	if hooks, _, err := g.Client.Repositories.ListHooks(g.ctxt, g.github_source.Organization, repo.Name, nil); err == nil {
		for _, hook := range hooks {
			if hook.GetName() != "web" {
				continue
			}
			if h, found := g.GetWebHook(hook, repo.WebHooks) ; found {
				h.identified = true
				repo.WebHooks[h.name] = h
				if h.Update(hook) {
					if _, _, err := g.Client.Repositories.EditHook(g.ctxt, g.github_source.Organization, repo.Name, hook.GetID(), hook); err != nil {
						log.Print(ret.Errorf("Failed to update '%s'. %s", h.name, err))
						return
					} else {
						log.Print(ret.StatusAdd("WebHook '%s' updated.", h.name))
					}
				}
			} else if repo.WebHookPolicy == "sync" || repo.WebHookPolicy == "" {
				if _, err := g.Client.Repositories.DeleteHook(g.ctxt, g.github_source.Organization, repo.Name, hook.GetID()); err != nil {
					log.Print(ret.Errorf("Failed to delete '%s'. %s", hook.Config["url"], err))
					return
				} else {
					log.Print(ret.StatusAdd("WebHook '%s' removed.", hook.Config["url"]))
				}
			} else {
				log.Print(ret.StatusAdd("WebHook '%s' not managed. Ignored. (policy = 'manage')", hook.Config["url"]))
			}

		}
	}
	for name, hook := range repo.WebHooks {
		if hook.identified {
			continue
		}

		Config := make(map[string]interface{})

		Config["url"] = hook.Url
		Config["insecure_ssl"] = g.SetWebHookInsecure(hook.SSLCheck)
		Config["content_type"] = hook.ContentType

		integ_name := "web" // webhook integration type in a dedicated webhook tab.
		new_hook := github.Hook{
			Name: &integ_name,
			Config: Config,
			Events: hook.Events,
		}

		hook.HookEnabled(&new_hook)

		if _, resp, err := g.Client.Repositories.CreateHook(g.ctxt, g.github_source.Organization, repo.Name, &new_hook); err != nil {
			log.Print(ret.Errorf("Failed to create '%s'. %s", name, err))
			return
		} else if resp.StatusCode != 201 {
			log.Print(ret.Errorf("Failed to create '%s'. %s", name, resp.Status))
			return
		}
		log.Print(ret.StatusAdd("WebHook '%s' created.", name))
		hook.identified = true
	}
	return true
}
