package main

import (
	"log"
	"strings"
	"fmt"
)

type RepositoryStruct  struct { // Used to stored the yaml source file. Not used to respond to the API requester.
	Name string                 // Name of the Repo
	Flow string                 // Flow applied on the repo.
	Description string          // Title in github repository
	Users map[string]string     // Collection of users role
	Groups map[string]string     // Collection of groups role
	// Following data are used at runtime but not saved. Used to respond to the API.
	exist bool                      // True if the repo exist.
	remotes map[string]string       // k: remote name, v: remote url
	branchConnect map[string]string // k: local branch name, v: remote/branch
}

func (r *RepositoryStruct)set(repo RepoAddStruct, remotes, branchConnect map[string]string) *RepositoryStruct {
	if r == nil {
		r = new(RepositoryStruct)
	}
	r.Name = repo.Name
	r.Description = repo.Title
	r.Flow = repo.Flow
	r.AddUsers(repo.Users)
	r.AddGroups(repo.Groups)
	r.remotes = remotes
	r.branchConnect = branchConnect
	return r
}

func (r *RepositoryStruct)AddUsers(users string) {
	if r.Users == nil {
		r.Users = make(map[string]string)
	}
	for _, user_role := range strings.Split(users, ",") {
		user_role_array := strings.Split(user_role, ":")
		user := user_role_array[0]
		role := user_role_array[1]
		if user == "" {
			log.Printf("Invalid user:role '%s' combination", user_role)
			continue
		}
		if role == "" {
			role = "developer"
			log.Printf("Role not defined for user '%s'. Using default 'developer'.", user)
		}
		r.Users[user] = role
	}
}

func (r *RepositoryStruct)AddGroups(groups string) {
	if r.Groups == nil {
		r.Groups = make(map[string]string)
	}
	for _, group_role := range strings.Split(groups, ",") {
		group_role_array := strings.Split(group_role, ":")
		group := group_role_array[0]
		role := group_role_array[1]
		if group == "" {
			log.Printf("Invalid group:role '%s' combination", group_role)
			continue
		}
		if role == "" {
			role = "developer"
			log.Printf("Role not defined for group '%s'. Using default 'developer'.", group)
		}
		r.Groups[group] = role
	}

}

func (r *RepositoryStruct)Update(repo RepoChangeStruct) (count int){
	if r.Description != repo.Title {
		r.Description = repo.Title
		count++
	}

	if r.Flow != repo.Flow {
		r.Flow = repo.Flow
		count++
	}

	// TODO: Be able to update the users/group list and their rights.

	return
}

// DoUpdateIn GitHubStruct with data from request.
//
func (r *RepoInstanceStruct) DoUpdateIn(g *GitHubStruct) (Updated bool, err, mess string) {
	if r.Add.Name != "" {
		// Add repo request type
		if g.AddRepo(r.Add.Name, r.Add) {
			Updated = true
			mess = fmt.Sprintf("New Repository '%s' added.", r.Add.Name)
		} else {
			err = fmt.Sprintf("Repository '%s' already exist.", r.Add.Name)
		}
	}

	if r.Change.Name != "" {
		// Change repo request type
		repo := g.github_source.Repos[r.Change.Name]
		if repo.Update(r.Change) > 0 {
			Updated = true
			mess = fmt.Sprintf("Repository '%s' updated.", r.Change.Name)
		} else {
			err = fmt.Sprintf("Repository '%s' doesn't exist. You must add it first.", r.Change.Name)
		}
	}
	return
}
