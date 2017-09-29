package main

import (
	"bytes"
	"fmt"
	"github.com/forj-oss/goforjj"
	"log"
	"os"
	"path"
)

// This file describes how we generate source from templates.

// loop on files to simply copy
func (p *JenkinsPlugin) copy_source_files(instance_name string, ret *goforjj.PluginData, status *bool) (_ error) {
	for file, desc := range p.sources {
		source_status := false
		src := path.Join(p.template_dir, desc.Source)
		dest := path.Join(p.source_path, desc.Source)
		parent := path.Dir(dest)

		if parent != "." {
			if _, err := os.Stat(parent); err != nil {
				log.Printf("Creating '%s'.", parent)
				if err = os.MkdirAll(parent, 0755); err != nil {
					log.Printf(ret.Errorf("Unable to copy '%s' to '%s'. %s.", src, dest, err))
					return
				}
			}
		}
		var dest_md5 []byte
		if m5, err := md5sum(dest); err == nil {
			dest_md5 = m5
		}
		if _, err, m5 := Copy(src, dest); err != nil {
			log.Printf(ret.Errorf("Unable to copy '%s' to '%s'. %s.", src, dest, err))
			return
		} else {
			if dest_md5 == nil || !bytes.Equal(dest_md5, m5) {
				IsUpdated(&source_status)
			}
		}

		if u, err := set_rights(dest, desc.Chmod); err != nil {
			ret.Errorf("%s", err)
			return err
		} else if u {
			IsUpdated(&source_status)
		}

		if source_status {
			IsUpdated(status)
			log.Printf("Copied '%s' to '%s'", src, dest)
			log.Printf(ret.StatusAdd("%s (%s) copied.", file, desc.Source))
			ret.AddFile(path.Join(instance_name, desc.Source))
		} else {
			log.Printf("'%s' not updated.", dest)
		}
	}
	return
}

func set_rights(file string, rights os.FileMode) (updated bool, _ error) {
	if rights == 0 {
		log.Printf("No rights to apply to %s.", file)
		return
	}

	log.Printf("Checking %s rights.", file)

	var rightsb os.FileMode
	stat_found := false
	if r, err := os.Stat(file); err == nil {
		rightsb = r.Mode()
		stat_found = true
	}
	if err := os.Chmod(file, rights); err != nil {
		return false, fmt.Errorf("Unable to set rights to '%s' with '%d'. %s", file, rights, err)
	}
	if stat_found {
		updated = (rightsb != rights)
	} else {
		updated = true
		log.Printf("%s rights updated from %d to %d.", file, rightsb, rights)
	}
	return
}

// loop on templates to use to generate source files
// The based data used for template is conform to the content of
// the forjj-jenkins.yaml file
// See YamlJenkins in jenkins_plugin.go
func (p *JenkinsPlugin) generate_source_files(instance_name string, ret *goforjj.PluginData, status *bool) (_ error) {
	for file, desc := range p.templates {
		if s, err := desc.Generate(p.yaml, p.template_dir, p.source_path, desc.Template); err != nil {
			log.Printf(ret.Errorf("%s", err))
			return err
		} else if s {
			ret.AddFile(path.Join(instance_name, desc.Template))
			log.Printf(ret.StatusAdd("%s (%s) generated", file, desc.Template))
			IsUpdated(status)
		} else {
			log.Printf("%s (%s) not updated", file, desc.Template)
		}
	}

	return
}
