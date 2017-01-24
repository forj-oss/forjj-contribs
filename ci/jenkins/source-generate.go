package main

import (
    "github.com/forj-oss/goforjj"
    "log"
    "path"
    "os"
    "fmt"
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

        if err := set_rights(dest, desc.Chmod) ; err != nil {
            ret.Errorf("%s", err)
            return
        }

        log.Printf(ret.StatusAdd("%s (%s) copied.", file, desc.Source))
        ret.AddFile(path.Join(instance_name, desc.Source))
    }
    return true
}

func set_rights(file string, rights os.FileMode) error {
    if rights != 0 {
        log.Printf("Setting %s rights to %d.", file, rights)
        os.Chmod(file, rights)
        if err := os.Chmod(file, rights) ; err != nil {
            return fmt.Errorf("Unable to set rights to '%s' with '%d'. %s", file, rights, err)
        }
    } else {
        log.Printf("No rights to apply to %s.", file)
    }
    return nil
}

// loop on templates to use to generate source files
// The based data used for template is conform to the content of
// the forjj-jenkins.yaml file
// See YamlJenkins in jenkins_plugin.go
func (p *JenkinsPlugin)generate_source_files(instance_name string, ret *goforjj.PluginData) (status bool) {
    for file, desc := range p.templates {
        if err := desc.Generate(p.yaml, p.template_dir, p.source_path, desc.Template) ; err != nil {
            log.Printf(ret.Errorf("%s", err))
            return
        }

        ret.AddFile(path.Join(instance_name, desc.Template))
        log.Printf(ret.StatusAdd("%s (%s) generated", file, desc.Template))
    }
    status = true
    return
}

