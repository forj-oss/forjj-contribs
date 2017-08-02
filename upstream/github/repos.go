package main

import (
	"fmt"
	"github.com/forj-oss/goforjj"
	"log"
	"strconv"
	"strings"
)

type RepositoryStruct struct { // Used to stored the yaml source file. Not used to respond to the API requester.
	Name         string            // Name of the Repo
	Flow         string            `yaml:",omitempty"`    // Flow applied on the repo.
	Description  string            `yaml:",omitempty"`    // Title in github repository
	Disabled     bool              `yaml:",omitempty"`    // disable the repository (became private with no team/collaborators)
	IssueTracker bool              `yaml:"issue_tracker"` // Issue tracker option
	Users        map[string]string // Collection of users role
	Groups       map[string]string // Collection of groups role
	// Following data are used at runtime but not saved. Used to respond to the API.
	exist         bool                                   // True if the repo exist.
	remotes       map[string]goforjj.PluginRepoRemoteUrl // k: remote name, v: remote urls
	branchConnect map[string]string                      // k: local branch name, v: remote/branch
}

func (r *RepositoryStruct) set(
	repo *RepoInstanceStruct,
	remotes map[string]goforjj.PluginRepoRemoteUrl,
	branchConnect map[string]string,
) *RepositoryStruct {
	if r == nil {
		r = new(RepositoryStruct)
	}
	r.Name = repo.Name
	r.Description = repo.Title
	if v, err := strconv.ParseBool(repo.Issue_tracker); err == nil {
		r.IssueTracker = v
		log.Printf("Issue_tracker '%s' => %t", repo.Issue_tracker, v)
	} else {
		log.Printf("IssueTracker has an invalid boolean string representation '%s'. Ignored. Tracker is set to true.",
			repo.Issue_tracker)
		r.IssueTracker = true
	}
	r.Flow = repo.Flow
	r.AddUsers(repo.Users)
	r.AddGroups(repo.Groups)
	r.remotes = remotes
	r.branchConnect = branchConnect
	return r
}

func (r *RepositoryStruct) AddUsers(users string) {
	if r.Users == nil {
		r.Users = make(map[string]string)
	}
	for _, user_role := range strings.Split(users, ",") {
		user_role_array := strings.Split(user_role, ":")
		user := ""
		role := ""
		if users_num := len(user_role_array); users_num >= 2 {
			user = user_role_array[0]
			role = user_role_array[1]
		} else {
			if users_num == 1 {
				user = user_role_array[0]
			}
		}
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

func (r *RepositoryStruct) AddGroups(groups string) {
	if r.Groups == nil {
		r.Groups = make(map[string]string)
	}
	for _, group_role := range strings.Split(groups, ",") {
		group_role_array := strings.Split(group_role, ":")
		group := ""
		role := ""
		if groups_num := len(group_role_array); groups_num >= 2 {
			group = group_role_array[0]
			role = group_role_array[1]
		} else {
			if groups_num == 1 {
				group = group_role_array[0]
			}
		}
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

func (r *RepositoryStruct) IsValid(repo_name string) (err error) {
	if r.Name == "" {
		err = fmt.Errorf("Invalid repository '%s'. Name is empty.", repo_name)
		return
	}
	if r.Name != repo_name {
		err = fmt.Errorf("Invalid repository '%s'. Name must be equal to '%s'. But the repo name is set to '%s'.",
			repo_name, repo_name, r.Name)
		return
	}
	return
}

func (r *RepoInstanceStruct) IsValid(repo_name string, ret *goforjj.PluginData) (valid bool) {
	if r.Name == "" {
		ret.Errorf("Invalid repository '%s'. Name is empty.", repo_name)
		return
	}
	if r.Name != repo_name {
		ret.Errorf("Invalid repository '%s'. Name must be equal to '%s'. But the repo name is set to '%s'.",
			repo_name, repo_name, r.Name)
		return
	}
	valid = true
	return
}
