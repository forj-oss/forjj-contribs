package main

import (
    "github.hpe.com/christophe-larsonneur/goforjj"
    "log"
    "path"
    "os"
    "text/template"
)

const jenkins_file = "forjj-jenkins.yaml"

// return true if instance doesn't exist.
func (r *CreateReq) check_source_existence(ret *goforjj.PluginData) (p *JenkinsPlugin, status bool) {
    if r.Groups.Source.Name == "" {
        log.Printf(ret.Errorf("Missing jenkins instance Name"))
        return
    }
    log.Printf("Checking Jenkins source code existence.")
    src := path.Join(r.ForjjSourceMount, r.Groups.Source.Name)
    if _, err := os.Stat(path.Join(src, jenkins_file)) ; err == nil {
        log.Printf(ret.Errorf("Unable to create the jenkins source code for instance name '%s' which already exist.\nUse update to update it (or update %s), and maintain to update github according to his configuration. %s.", src, src, err))
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

    log.Printf(ret.StatusAdd("environment checked."))
    return p, true
}

func (r *JenkinsPlugin)create_jenkins_sources(ret *goforjj.PluginData) (status bool) {

    if err := r.LoadTemplatesDef() ; err != nil {
        log.Printf(ret.Errorf("%s", err))
        return
    }

    if err := r.DefineSources() ; err != nil {
        log.Printf(ret.Errorf("%s", err))
        return
    }

    if ! r.copy_source_files(ret) {
        return
    }

    if ! r.generate_source_files(ret) {
        return
    }

    return true
}

// loop on files to simply copy
func (p *JenkinsPlugin)copy_source_files(ret *goforjj.PluginData) (status bool) {
    for file, desc := range p.sources {
        src := path.Join(p.template_dir, desc.Source)
        dest := path.Join(p.source_path, desc.Source)
        parent := path.Dir(dest)

        if  parent != "." {
            if _, err := os.Stat(parent) ; err != nil {
                os.MkdirAll(parent, 0755)
            }
        }
        if _, err := Copy(src, path.Join(p.source_path, desc.Source)) ; err != nil {
            log.Printf(ret.Errorf("Unable to copy '%s' to '%s'. %s.", src, dest, err))
            return
        }
        log.Printf(ret.StatusAdd("%s (%s) copied", file, desc.Source))
    }
    return true
}

// loop on templates to use to generate source files
func (p *JenkinsPlugin)generate_source_files(ret *goforjj.PluginData) (status bool) {

    templates := make([]string, 10)

    for _, desc := range p.templates {
        templates = append(templates, path.Join(p.template_dir, desc.Template))
    }

    t, err := template.New("jenkins").Funcs(template.FuncMap{}).ParseFiles(templates...)

    if err != nil {
        log.Printf(ret.Errorf("Template issue: %s", err))
        return
    }

    for file, desc := range p.templates {
        dest := path.Join(p.source_path, desc.Template)
        parent := path.Dir(dest)

        if  parent != "." {
            if _, err := os.Stat(parent) ; err != nil {
                os.MkdirAll(parent, 0755)
            }
        }

        if out, err := os.Create(dest) ; err != nil {
            log.Printf(ret.Errorf("Unable to create %s. %s.", dest, err))
            return
        } else {
            t.ExecuteTemplate(out, desc.Template, p.data)
            out.Close()
        }

        log.Printf(ret.StatusAdd("%s (%s) generated", file, desc.Template))
    }
    return true
}
