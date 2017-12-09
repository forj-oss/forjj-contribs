// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"github.com/forj-oss/goforjj"
	"log"
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
						log.Printf("Failed to update '%s'. %s", hook.GetName(), err)
					}
				}
			}

		}
	}
	for name, hook := range g.github_source.WebHooks {

	}
	return true
}
