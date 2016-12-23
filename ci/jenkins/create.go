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

	ret.CommitMessage = "Creating initial jenkins source files."
	return true
}

func (r *CreateArgReq) SaveMaintainOptions(ret *goforjj.PluginData) {
	if ret.Options == nil {
		ret.Options = make(map[string]goforjj.PluginOption)
	}
}
