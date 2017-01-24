package main

import "fmt"

func (t *GithubStruct)SetFrom(d *GithubStruct) bool {
	if t == nil {
		return false
	}
	SetIfSet(&t.ApiUrl, d.ApiUrl)
	SetIfSet(&t.Repo, d.Repo)
	SetIfSet(&t.RepoOwner, d.RepoOwner)
	return true
}

func (t *GithubStruct)UpdateFrom(d *GithubStruct) bool {
	if t == nil {
		return false
	}
	SetIfSet(&t.ApiUrl, d.ApiUrl)
	SetIfSet(&t.Repo, d.Repo)
	SetIfSet(&t.RepoOwner, d.RepoOwner)
	return true
}

func (t *GithubStruct)GetUpstream() string {
	return fmt.Sprintf("%s/%s/%s", t.ApiUrl, t.RepoOwner, t.Repo)
}
