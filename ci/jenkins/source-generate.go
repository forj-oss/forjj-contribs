package main

import (
    "github.hpe.com/christophe-larsonneur/goforjj"
    "log"
    "path"
    "os"
    "text/template"
    "io/ioutil"
)

// This file describes how we generate source from templates.

// loop on files to simply copy
func (p *JenkinsPlugin)copy_source_files(instance_name string, ret *goforjj.PluginData) (status bool) {
    for file, desc := range p.sources {
        src := path.Join(p.template_dir, desc.Source)
        dest := path.Join(p.source_path, desc.Source)
        parent := path.Dir(dest)

        if _, err := os.Stat(dest) ; err == nil {
            log.Printf(ret.StatusAdd("%s (%s) already exist. Not copied.", file, desc.Source))
            ret.AddFile(path.Join(instance_name, desc.Source))
            continue
        }
        log.Printf("Copying '%s' to '%s'", src, dest)

        if  parent != "." {
            if _, err := os.Stat(parent) ; err != nil {
                log.Printf("Creating '%s'.", parent)
                os.MkdirAll(parent, 0755)
            }
        }
        if _, err := Copy(src, path.Join(p.source_path, desc.Source)) ; err != nil {
            log.Printf(ret.Errorf("Unable to copy '%s' to '%s'. %s.", src, dest, err))
            return
        }
        log.Printf(ret.StatusAdd("%s (%s) copied.", file, desc.Source))
        ret.AddFile(path.Join(instance_name, desc.Source))
    }
    return true
}

// loop on templates to use to generate source files
// The based data used for template is conform to the conent of
// the forjj-jenkins.yaml file
// See YamlJenkins in jenkins_plugin.go
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
