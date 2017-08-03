// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"github.com/forj-oss/goforjj"
	"log"
	"os"
	"path"
)

// Return ok if the jenkins instance exist
func (r *UpdateReq) check_source_existence(ret *goforjj.PluginData) (p *JenkinsPlugin, status bool) {
	log.Print("Checking Jenkins source code existence.")
	if _, err := os.Stat(r.Forj.ForjjSourceMount); err != nil {
		ret.Errorf("Unable to update jenkins instances. '%s' is inexistent or innacessible. %s", r.Forj.ForjjSourceMount, err)
		return
	}

	src_path := path.Join(r.Forj.ForjjSourceMount, r.Forj.ForjjInstanceName)

	p = new_plugin(src_path)

	p.template_dir = *cliApp.params.template_dir
	templatef := path.Join(p.template_dir, template_file)
	if _, err := os.Stat(templatef); err != nil {
		log.Printf(ret.Errorf("Unable to find templates definition file '%s'. %s.", templatef, err))
		return
	}

	p.template_file = templatef
	ret.StatusAdd("environment checked.")
	return p, true
}

func (r *JenkinsPlugin) update_jenkins_sources(instance_name string, ret *goforjj.PluginData) (status bool) {
	if err := r.DefineSources(); err != nil {
		log.Printf(ret.Errorf("%s", err))
		return
	}

	log.Print("Start copying source files...")
	status = r.copy_source_files(instance_name, ret)

	log.Print("Start Generating source files...")
	status = r.generate_source_files(instance_name, ret) || status

	status = r.generate_jobsdsl(instance_name, ret)  || status

	if status {
		ret.CommitMessage = "Creating initial jenkins source files."
	} else {
		ret.StatusAdd("No update detected.")
	}
	return
}

// Function which adds maintain options as part of the plugin answer in create/update phase.
// forjj won't add any driver name because 'maintain' phase read the list of drivers to use from forjj-maintain.yml
// So --git-us is not available for forjj maintain.
func (r *UpdateArgReq) SaveMaintainOptions(ret *goforjj.PluginData) {
	if ret.Options == nil {
		ret.Options = make(map[string]goforjj.PluginOption)
	}
}

func addMaintainOptionValue(options map[string]goforjj.PluginOption, option, value, defaultv, help string) goforjj.PluginOption {
	opt, ok := options[option]
	if ok && value != "" {
		opt.Value = value
		return opt
	}
	if !ok {
		opt = goforjj.PluginOption{Help: help}
		if value == "" {
			opt.Value = defaultv
		} else {
			opt.Value = value
		}
	}
	return opt
}
