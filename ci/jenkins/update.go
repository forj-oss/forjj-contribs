package main

import (
    "github.hpe.com/christophe-larsonneur/goforjj"
    "log"
    "path"
    "os"
)

// Return ok if the jenkins instance exist
func (r *UpdateReq) check_source_existence(ret *goforjj.PluginData) (p *JenkinsPlugin, status bool) {
    if r.Groups.Source.Name == "" {
        ret.Errorf("Missing jenkins instance Name")
        return
    }
    log.Printf("Checking Jenkins source code existence.")
    src := path.Join(r.ForjjSourceMount, r.Groups.Source.Name)
    if _, err := os.Stat(path.Join(src, jenkins_file)) ; err == nil {
        ret.Errorf("Unable to create the jenkins source code for instance name '%s' which already exist.\nUse update to update it (or update %s), and maintain to update github according to his configuration.", src)
        return
    }

    p = r.Groups.new_plugin(src)

    p.template_dir = *cliApp.params.template_dir
    templatef := path.Join(p.template_dir, template_file)
    if _, err := os.Stat(templatef) ; err == nil {
        log.Printf(ret.Errorf("Unable to find templates definition file '%s'. %s.", templatef, err))
        return
    }

    p.template_file = templatef

    ret.StatusAdd("environment checked.")
    return p, true
}

func (r *JenkinsPlugin)update_jenkins_sources(ret *goforjj.PluginData) (status bool) {
    return true
}
