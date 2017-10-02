// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"github.com/forj-oss/goforjj"
	"log"

	"os"
	"path"
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

	log.Printf(ret.StatusAdd("environment checked."))
	return
}

// We assume template source file is loaded.
func (r *JenkinsPlugin) create_jenkins_sources(instance_name string, ret *goforjj.PluginData) (err error) {

	if err = r.DefineSources(); err != nil {
		log.Printf(ret.Errorf("%s", err))
		return err
	}

	log.Print("Start copying source files...")
	if err = r.copy_source_files(instance_name, ret, nil); err != nil {
		return
	}

	log.Print("Start Generating source files...")
	if err = r.generate_source_files(instance_name, ret, nil); err != nil {
		return
	}

	if err = r.generate_jobsdsl(instance_name, ret, nil); err != nil {
		return
	}

	return
}

// add_projects add project data in the jenkins.yaml file
func (jp *JenkinsPlugin) add_projects(req *CreateReq, ret *goforjj.PluginData) error {
	projects := ProjectInfo{}
	projects.set_project_info(req.Forj.ForjCommonStruct)
	projects.set_infra_remote(req.Objects.App[req.Forj.ForjjInstanceName].SeedJobRepo)
	return projects.set_projects_to(req.Objects.Projects, jp, ret, nil)
}

// generate_jobsdsl generate any missing job-dsl source file.
// TODO: Support for different Repository path that Forjj have to checkout appropriately.
func (p *JenkinsPlugin) generate_jobsdsl(instance_name string, ret *goforjj.PluginData, status *bool) (err error) {
	if p.yaml.Projects == nil {
		return
	}
	if err = p.yaml.Projects.Generates(p, instance_name, ret, status); err != nil {
		log.Print(ret.Errorf("%s", err))
	}
	return
}

func (r *CreateArgReq) SaveMaintainOptions(ret *goforjj.PluginData) {
	if ret.Options == nil {
		ret.Options = make(map[string]goforjj.PluginOption)
	}
}
