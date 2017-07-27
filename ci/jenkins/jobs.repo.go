package main

import (
	"fmt"
	"github.com/forj-oss/goforjj"
	"log"
	"os"
	"path"
)

type Projects struct {
	DslRepo    string
	DslPath    string
	InfraPath  string
	infra_repo bool
	List       map[string]Project
}

type Project struct {
	Name       string
	SourceType string
	Github     GithubStruct `yaml:",omitempty"`
	Git        GitStruct    `yaml:",omitempty"`
	all        *Projects
}

func NewProjects(InstanceName, repo, Dslpath string, infra_repo bool) *Projects {
	p := new(Projects)
	p.DslPath = Dslpath
	p.DslRepo = repo
	if infra_repo {
		p.InfraPath = path.Join("apps", "ci", InstanceName)
		p.DslPath = "jobs-dsl"
	}

	p.infra_repo = infra_repo
	p.List = make(map[string]Project)
	return p
}

func (p *Projects) AddGithub(name string, d *GithubStruct) bool {
	data := new(GithubStruct)
	data.SetFrom(d)
	p.List[name] = Project{Name: name, SourceType: "github", Github: *data, all: p}
	return true
}

func (p *Projects) AddGit(name string, d *GitStruct) bool {
	data := new(GitStruct)
	data.SetFrom(d)
	p.List[name] = Project{Name: name, SourceType: "git", Git: *data, all: p}
	return true
}

func (p *Project) Remove() bool {
	return true
}

func (p *Project) Add() error {
	return nil
}

// Generates Jobs-dsl files in the given checked-out GIT repository.
func (p *Projects) Generates(instance_name, template_dir, repo_path string, ret *goforjj.PluginData) (bool, error) {
	if f, err := os.Stat(repo_path); err != nil {
		return false, err
	} else {
		if !f.IsDir() {
			return false, fmt.Errorf("Repo cloned path '%s' is not a directory.", repo_path)
		}
	}

	jobs_dsl_path := path.Join(repo_path, p.DslPath)
	if f, err := os.Stat(jobs_dsl_path); err != nil {
		os.MkdirAll(jobs_dsl_path, 0755)
	} else {
		if !f.IsDir() {
			return false, fmt.Errorf("path '%s' is not a directory.", repo_path)
		}
	}

	tmpl := new(TmplSource)
	tmpl.Template = "jobs-dsl/job-dsl.tmpl"
	tmpl.Chmod = 0644

	for name, prj := range p.List {
		if err := tmpl.Generate(prj, template_dir, jobs_dsl_path, name+".groovy"); err != nil {
			log.Printf("Unable to generate '%s'. %s",
				path.Join(jobs_dsl_path, name+".groovy"), ret.Errorf("%s", err))
			continue
		}
		ret.AddFile(path.Join(instance_name, jobs_dsl_path, name+".groovy"))
		log.Printf(ret.StatusAdd("Project '%s' (%s) generated", name, path.Join(p.DslPath, name+".groovy")))
	}
	return true, nil
}
