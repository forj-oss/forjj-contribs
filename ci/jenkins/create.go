package main

import (
    "github.hpe.com/christophe-larsonneur/goforjj"
    "log"
    "path"
    "os"
    "text/template"
    "io/ioutil"
)

// return true if instance doesn't exist.
func (r *CreateReq) check_source_existence(ret *goforjj.PluginData) (p *JenkinsPlugin, httpCode int) {
    log.Printf("Checking Jenkins source code existence.")
    src := path.Join(r.ForjjSourceMount, r.ForjjInstanceName)
    if _, err := os.Stat(path.Join(src, jenkins_file)) ; err == nil {
        log.Printf(ret.Errorf("Unable to create the jenkins source code for instance name '%s' which already exist.\nUse 'update' to update it (or update %s), and 'maintain' to update jenkins according to his configuration.", r.ForjjInstanceName, src))
        return nil, 419 // Abort message returned to forjj.
    }

    p = new_plugin(src)

    p.template_dir = *cliApp.params.template_dir
    templatef := path.Join(p.template_dir, template_file)
    if _, err := os.Stat(templatef) ; err != nil {
        log.Printf(ret.Errorf("Unable to find templates definition file '%s'. %s.", templatef, err))
        return
    }

    p.template_file = templatef

    log.Printf(ret.StatusAdd("environment checked."))
    return
}

func (r *JenkinsPlugin)create_jenkins_sources(instance_name string, ret *goforjj.PluginData) (status bool) {

    if err := r.LoadTemplatesDef() ; err != nil {
        log.Printf(ret.Errorf("%s", err))
        return
    }

    if err := r.DefineSources() ; err != nil {
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

// loop on files to simply copy
func (p *JenkinsPlugin)copy_source_files(instance_name string, ret *goforjj.PluginData) (status bool) {
    for file, desc := range p.sources {
        src := path.Join(p.template_dir, desc.Source)
        dest := path.Join(p.source_path, desc.Source)
        parent := path.Dir(dest)

        log.Printf("Copying '%s' to '%s'", src, dest)

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
        ret.AddFile(path.Join(instance_name, desc.Source))
    }
    return true
}

// loop on templates to use to generate source files
func (p *JenkinsPlugin)generate_source_files(instance_name string, ret *goforjj.PluginData) (status bool) {
    for file, desc := range p.templates {
        src := path.Join(p.template_dir, desc.Template)
        dest := path.Join(p.source_path, desc.Template)
        parent := path.Dir(dest)

        if  parent != "." {
            if _, err := os.Stat(parent) ; err != nil {
                os.MkdirAll(parent, 0755)
            }
        }

        var data string
        if b, err := ioutil.ReadFile(src) ; err != nil {
            log.Printf(ret.Errorf("Load issue. %s", err))
        } else {
            data = string(b)
        }

        t, err := template.New(src).Funcs(template.FuncMap{}).Parse(data)
        if err != nil {
            log.Printf(ret.Errorf("Template issue. %s", err))
            return
        }

        if out, err := os.Create(dest) ; err != nil {
            log.Printf(ret.Errorf("Unable to create %s. %s.", dest, err))
            return
        } else {
            if err := t.Execute(out, p.yaml) ; err != nil {
                log.Printf(ret.Errorf("Unable to interpret %s. %s.", dest, err))
            }
            out.Close()
        }
        ret.AddFile(path.Join(instance_name, desc.Template))
        log.Printf(ret.StatusAdd("%s (%s) generated", file, desc.Template))
    }
    return true
}
