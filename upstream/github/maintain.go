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

			if  hook.Name == nil {
				continue
			}
			if h, found := g.github_source.WebHooks[*hook.Name] ; found {
				h.identified = true
				g.github_source.WebHooks[*hook.Name] = h
				if h.Update(hook) {
					if _, _, err := g.Client.Organizations.EditHook(g.ctxt, g.github_source.Organization, hook.GetID(), hook) ; err != nil {
						log.Print(ret.Errorf("Failed to update '%s'. %s", hook.GetName(), err))
						return
					}
				}
			} else if  g.github_source.WebHookPolicy == "sync" {
				if _, err := g.Client.Organizations.DeleteHook(g.ctxt, g.github_source.Organization, hook.GetID()); err != nil {
					log.Print(ret.Errorf("Failed to delete '%s'. %s", hook.GetName(), err))
					return
				} else {
					log.Print(ret.StatusAdd("Org Hook '%s' removed.", hook.GetName()))
				}
			} else {
				log.Print(ret.StatusAdd("Org Hook '%s' not managed. (org webhook policy = 'manage')", hook.GetName()))
			}
		}
	}
	for name, hook := range g.github_source.WebHooks {
		if hook.identified {
			continue
		}

		var Config map[string]interface{}

		Config["url"] = hook.Url
		Config["insecure_ssl"] = hook.SSLCheck
		Config["content_type"] = hook.ContentType


		new_hook := github.Hook{
			Name: &name,
			Config: Config,
			Events: hook.Events,
		}

		hook.HookEnabled(&new_hook)

		if _, _, err := g.Client.Organizations.CreateHook(g.ctxt, g.github_source.Organization, &new_hook) ; err != nil {
			log.Print(ret.Errorf("Failed to delete '%s'. %s", name, err))
		}
		log.Print(ret.StatusAdd("Org Hook '%s created.", new_hook.GetName()))
		hook.identified = true
	}
	return true
}

func (g *GitHubStruct) MaintainHooks(repo *RepositoryStruct, ret *goforjj.PluginData) (_ bool) {
	// organization level
	if hooks, _, err := g.Client.Repositories.ListHooks(g.ctxt, g.github_source.Organization, repo.Name, nil); err == nil {
		for _, hook := range hooks {

			if hook.Name == nil {
				continue
			}
			if h, found := repo.WebHooks[*hook.Name]; found {
				h.identified = true
				repo.WebHooks[*hook.Name] = h
				if h.Update(hook) {
					if _, _, err := g.Client.Repositories.EditHook(g.ctxt, g.github_source.Organization, repo.Name, hook.GetID(), hook); err != nil {
						log.Print(ret.Errorf("Failed to update '%s'. %s", hook.GetName(), err))
						return
					} else {
						log.Print(ret.StatusAdd("Hook '%s updated.", hook.GetName()))
					}
				}
			} else if repo.WebHookPolicy == "sync" {
				if _, err := g.Client.Repositories.DeleteHook(g.ctxt, g.github_source.Organization, repo.Name, hook.GetID()); err != nil {
					log.Print(ret.Errorf("Failed to delete '%s'. %s", hook.GetName(), err))
					return
				} else {
					log.Print(ret.StatusAdd("Hook '%s' removed.", hook.GetName()))
				}
			} else {
				log.Print(ret.StatusAdd("Hook '%s' not managed. Ignored. (policy = 'manage')", hook.GetName()))
			}

		}
	}
	for name, hook := range g.github_source.WebHooks {
		if hook.identified {
			continue
		}

		var Config map[string]interface{}

		Config["url"] = hook.Url
		Config["insecure_ssl"] = hook.SSLCheck
		Config["content_type"] = hook.ContentType

		new_hook := github.Hook{
			Name: &name,
			Config: Config,
			Events: hook.Events,
		}

		hook.HookEnabled(&new_hook)

		if _, _, err := g.Client.Repositories.CreateHook(g.ctxt, g.github_source.Organization, repo.Name, &new_hook); err != nil {
			log.Print(ret.Errorf("Failed to delete '%s'. %s", name, err))
			return
		}
		log.Print(ret.StatusAdd("Hook '%s created.", new_hook.GetName()))
		hook.identified = true
	}
	return true
}
