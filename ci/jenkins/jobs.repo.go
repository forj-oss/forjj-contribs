package main

import (
	"os"
	"fmt"
	"path"
	"log"
	"github.com/forj-oss/goforjj"
)

type Projects struct {
	DslRepo string
	DslPath string
	List map[string]Project
}

type Project struct {
	Name string
	Github *GithubStruct
	Git *GitStruct
	All *Projects
}

func NewProjects(repo, path string) *Projects {
	p := new(Projects)
	p.DslPath = path
	p.DslRepo = repo
	p.List = make(map[string]Project)
	return p
}

func (p *Projects)AddGithub(name string, d *GithubStruct) bool {
	data := new(GithubStruct)
	data.SetFrom(d)
	p.List[name] = Project{Name: name, Github: data, All: p}
	return true
}

func (p *Projects)AddGit(name string, d *GitStruct) bool {
	data := new(GitStruct)
	data.SetFrom(d)
	p.List[name] = Project{Name: name, Git: data, All: p}
	return true
}

func (p *Project)Remove() bool {
	return true
}

func (p* Project)Add() error {
	return nil
}

// Generates Jobs-dsl files in the given checked-out GIT repository.
func (p *Projects)Generates(instance_name, template_dir, repo_path string, ret *goforjj.PluginData) (bool, error) {
	if f, err := os.Stat(repo_path) ; err != nil {
		return false, err
	} else {
		if ! f.IsDir() {
			return false, fmt.Errorf("Repo cloned path '%s' is not a directory.", repo_path)
		}
	}

	jobs_dsl_path := path.Join(repo_path, p.DslPath)
	if f, err := os.Stat(jobs_dsl_path) ; err != nil {
		os.MkdirAll(jobs_dsl_path, 0755)
	} else {
		if ! f.IsDir() {
			return false, fmt.Errorf("path '%s' is not a directory.", repo_path)
		}
	}
	return true, nil

	tmpl := new(TmplSource)
	tmpl.Template = "jobs-dsl/job-dsl.tmpl"
	tmpl.Chmod = 0644

	for name, prj := range p.List {
		tmpl.Generate(prj, template_dir, jobs_dsl_path, name + ".groovy")
		ret.AddFile(path.Join(instance_name, jobs_dsl_path, name + ".groovy"))
		log.Printf(ret.StatusAdd("%s generated", path.Join(p.DslRepo, jobs_dsl_path, name + ".groovy")))
	}
	return true, nil
}
