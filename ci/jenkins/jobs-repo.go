package main

import (
	"fmt"
	"github.com/forj-oss/goforjj"
	"log"
	"os"
	"path"
	"strings"
)

type Projects struct {
	DslRepo    string
	DslPath    string
	InfraPath  string
	infra_repo bool
	infra_name string
	List       map[string]Project
}

type Project struct {
	Name       string
	SourceType string
	Github     GithubStruct `yaml:",omitempty"`
	Git        GitStruct    `yaml:",omitempty"`
	InfraRepo  bool         `yaml:",omitempty"`
	all        *Projects
}

type ProjectModel struct {
	Project Project
	Source  YamlJenkins
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

func (p *Projects) AddGithub(name string, d *GithubStruct, isInfra bool) bool {
	data := new(GithubStruct)
	data.SetFrom(d)
	p.List[name] = Project{Name: name, SourceType: "github", Github: *data, all: p, InfraRepo: isInfra}
	return true
}

func (p *Projects) AddGit(name string, d *GitStruct, isInfra bool) bool {
	data := new(GitStruct)
	data.SetFrom(d)
	p.List[name] = Project{Name: name, SourceType: "git", Git: *data, all: p, InfraRepo: isInfra}
	return true
}

func (p *Project) Remove() bool {
	return true
}

func (p *Project) Model(jp *JenkinsPlugin) (ret *ProjectModel) {
	ret = new(ProjectModel)
	ret.Project = *p
	ret.Source = jp.yaml
	return
}

func (p *Project) Add() error {
	return nil
}

// Generates Jobs-dsl files in the given checked-out GIT repository.
func (p *Projects) Generates(jp *JenkinsPlugin, instance_name string, ret *goforjj.PluginData, status *bool) (_ error) {
	template_dir := jp.template_dir
	repo_path := jp.source_path

	if f, err := os.Stat(repo_path); err != nil {
		return err
	} else {
		if !f.IsDir() {
			return fmt.Errorf(ret.Errorf("Repo cloned path '%s' is not a directory.", repo_path))
		}
	}

	jobs_dsl_path := path.Join(repo_path, p.DslPath)
	if f, err := os.Stat(jobs_dsl_path); err != nil {
		if err := os.MkdirAll(jobs_dsl_path, 0755); err != nil {
			return err
		}
	} else {
		if !f.IsDir() {
			return fmt.Errorf(ret.Errorf("path '%s' is not a directory.", repo_path))
		}
	}

	tmpl := new(TmplSource)
	tmpl.Template = "jobs-dsl/job-dsl.groovy"
	tmpl.Chmod = 0644

	for name, prj := range p.List {
		name = strings.Replace(name, "-", "_", -1)
		if u, err := tmpl.Generate(prj.Model(jp), template_dir, jobs_dsl_path, name+".groovy"); err != nil {
			log.Printf("Unable to generate '%s'. %s",
				path.Join(jobs_dsl_path, name+".groovy"), ret.Errorf("%s", err))
			return err
		} else if u {
			IsUpdated(status)
			ret.AddFile(path.Join(instance_name, jobs_dsl_path, name+".groovy"))
			log.Printf(ret.StatusAdd("Project '%s' (%s) generated", name, path.Join(p.DslPath, name+".groovy")))
		}
	}
	return nil
}
