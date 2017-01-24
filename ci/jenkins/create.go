// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"github.com/forj-oss/goforjj"
	"log"
	"os"
	"path"
	"net/url"
	"regexp"
	"fmt"
)

// return true if instance doesn't exist.
func (r *CreateReq) check_source_existence(ret *goforjj.PluginData) (p *JenkinsPlugin, httpCode int) {
	log.Printf("Checking Jenkins source code existence.")
	src := path.Join(r.Forj.ForjjSourceMount, r.Forj.ForjjInstanceName)
	if _, err := os.Stat(path.Join(src, jenkins_file)); err == nil {
		log.Printf(ret.Errorf("Unable to create the jenkins source code for instance name '%s' which already exist.\nUse 'update' to update it (or update %s), and 'maintain' to update jenkins according to his configuration.", r.Forj.ForjjInstanceName, src))
		return nil, 419 // Abort message returned to forjj.
	}

	p = new_plugin(src)

	p.template_dir = *cliApp.params.template_dir
	templatef := path.Join(p.template_dir, template_file)
	if _, err := os.Stat(templatef); err != nil {
		log.Printf(ret.Errorf("Unable to find templates definition file '%s'. %s.", templatef, err))
		return
	}

	p.template_file = templatef

	log.Printf(ret.StatusAdd("environment checked."))
	return
}

// We assume template source file is loaded.
func (r *JenkinsPlugin)create_jenkins_sources(instance_name string, ret *goforjj.PluginData) (status bool) {

	if err := r.DefineSources(); err != nil {
		log.Printf(ret.Errorf("%s", err))
		return
	}

	log.Print("Start copying source files...")
	if ! r.copy_source_files(instance_name, ret) {
		return
	}

	log.Print("Start Generating source files...")
	if ! r.generate_source_files(instance_name, ret) {
		return
	}

	if ! r.generate_jobsdsl(instance_name, ret) {
		return
	}

	ret.CommitMessage = "Creating initial jenkins source files."
	return true
}

// add_projects add project data in the jenkins.yaml file
func (r *JenkinsPlugin)add_projects(req *CreateReq, ret *goforjj.PluginData) (status bool) {
	if req.Forj.ForjjInfraUpstream == "" {
		ret.StatusAdd("Unable to add a new project without a remote GIT repository. Jenkins JobDSL requirement. " +
			"To enable this feature, add a remote GIT to your infra --infra-upstream or define the JobDSL Repository to clone.")
		return true
	}

	infra_remote := req.Forj.ForjjInfraUpstream
	ssh_format, _ := regexp.Compile(`^([a-z]+@)?(([a-z.-]+):)(/?\w[\w./-]*)?$`)
	if r := ssh_format.FindStringSubmatch(infra_remote) ; r != nil {
		infra_remote = fmt.Sprintf("ssh://%s%s/%s", r[1], r[2], r[4])
	}

	if v, err := url.Parse(infra_remote) ; err != nil {
		ret.Errorf("Infra remote URL issue. %s", err)
		return false
	} else {
		if v.Scheme == "" {
			ret.Errorf("Invalid Remote GIT repository Url '%s'. A scheme must exist.", infra_remote)
		}
	}
	// Initialize JobDSL structure
	r.yaml.Projects = NewProjects(infra_remote, "jobs-dsl")

	// Retrieve list of Repository (projects) to manage
	for name, prj := range req.Objects.Projects {
		switch prj.RemoteType {
		case "github":
			r.yaml.Projects.AddGithub(name, &prj.GithubStruct)
		case "git":
			r.yaml.Projects.AddGit(name, &prj.GitStruct)
		}
	}
	status = true
	return
}

// generate_jobsdsl generate any missing job-dsl source file.
// TODO: Support for different Repository path that Forjj have to checkout appropriately.
func (p *JenkinsPlugin)generate_jobsdsl(instance_name string, ret *goforjj.PluginData)(status bool) {
	if p.yaml.Projects == nil {
		return true // Nothing to do. But it is acceptable as not CORE.
	}
	if ok, err := p.yaml.Projects.Generates(instance_name, p.template_dir, p.source_path, ret) ; err != nil {
		log.Print(ret.Errorf("%s", err))
	} else {
		status = ok
	}
	return
}

func (r *CreateArgReq) SaveMaintainOptions(ret *goforjj.PluginData) {
	if ret.Options == nil {
		ret.Options = make(map[string]goforjj.PluginOption)
	}
}
