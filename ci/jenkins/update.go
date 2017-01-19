package main

import (
    "github.com/forj-oss/goforjj"
    "log"
    "path"
    "os"
)

// Return ok if the jenkins instance exist
func (r *UpdateReq) check_source_existence(ret *goforjj.PluginData) (p *JenkinsPlugin, status bool) {
    log.Printf("Checking Jenkins source code existence.")
    src := path.Join(r.Args.ForjjSourceMount, r.Args.ForjjInstanceName)
    if _, err := os.Stat(path.Join(src, jenkins_file)) ; err != nil {
        ret.Errorf("Unable to update the jenkins source code for instance name '%s' which doesn't exist.\nUse create to create it.", r.Args.ForjjInstanceName)
        return
    }

    p = new_plugin(src)

    p.template_dir = *cliApp.params.template_dir
    templatef := path.Join(p.template_dir, template_file)
    if _, err := os.Stat(templatef) ; err != nil {
        log.Printf(ret.Errorf("Unable to find templates definition file '%s'. %s.", templatef, err))
        return
    }

    p.template_file = templatef

    ret.StatusAdd("environment checked.")


    return p, true
}

// We assume template file were loaded.
func (r *JenkinsPlugin)update_jenkins_sources(instance_name string, ret *goforjj.PluginData) (status bool) {
    if err := r.DefineSources() ; err != nil {
        log.Printf(ret.Errorf("%s", err))
        return
    }

    log.Print("Start copying NEW source files...")
    if ! r.copy_source_files(instance_name, ret) {
        return
    }

    log.Print("Start re-generating source files...")
    if ! r.generate_source_files(instance_name, ret) {
        return
    }

    // Default commit message. Usually, forjj update has an end user commit message to apply instead. Or no commit to do at all.
    ret.CommitMessage = "Jenkins source files updated."

    return true
}
