package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
)

const template_file = "templates.yaml"

// Contains functions to manage source code from templates

type YamlTemplates struct {
	Defaults DefaultsStruct
	Features TmplFeatures
	Sources  TmplSources
	Run      map[string]RunStruct `yaml:"run_deploy"`
	Variants map[string]string
}

type DefaultsStruct struct {
	Dockerfile   DockerfileStruct
	JenkinsImage FinalImageStruct
}

type TmplFeatures struct {
	Common TmplFeaturesStruct
	Deploy map[string]TmplFeaturesStruct
}

type TmplFeaturesStruct []string

type TmplSources struct {
	Common TmplSourcesStruct
	Deploy map[string]TmplSourcesStruct
}

type TmplSourcesStruct map[string]TmplSource

type TmplSource struct {
	Chmod    os.FileMode
	Template string
	Source   string
	If       string `yaml:"if"` // If `If` is empty, the file will be ignored. otherwise the file will copied/generated
								// as usual.
}

type EnvStruct struct {
	Value string
	If string
}

type RunStruct struct {
	RunCommand string `yaml:"run"`
	Env map[string]EnvStruct `yaml:"environment"`
}

// Model creates the Model data used by gotemplates.
// The model is not updated until call to CleanModel()
func (p *JenkinsPlugin) Model() *JenkinsPluginModel {
	if JP_Model != nil {
		return JP_Model
	}
	JP_Model = new(JenkinsPluginModel)
	JP_Model.Source = p.yaml
	return JP_Model
}

func (p *JenkinsPlugin) CleanModel() {
	JP_Model = nil
}

func Evaluate(value string, data interface{}) (string, error) {
	var doc bytes.Buffer
	tmpl := template.New("jenkins_plugin_data")

	if !strings.Contains(value, "{{") {
		return value, nil
	}
	if _, err := tmpl.Parse(value); err != nil {
		return "", err
	}
	if err := tmpl.Execute(&doc, data); err != nil {
		return "", err
	}
	ret := doc.String()
	log.Printf("'%s' were interpreted to '%s'", value, ret)
	return ret, nil
}

//Load templates definition file from template dir.
func (p *JenkinsPlugin) LoadTemplatesDef() error {
	if d, err := ioutil.ReadFile(p.template_file); err != nil {
		return fmt.Errorf("Unable to load '%s'. %s.", p.template_file, err)
	} else {
		if err := yaml.Unmarshal(d, &p.templates_def); err != nil {
			return fmt.Errorf("Unable to load yaml file format '%s'. %s.", p.template_file, err)
		}
	}
	return nil
}

// Load list of files to copy and files to generate
func (p *JenkinsPlugin) DefineSources() error {
	// load all features
	p.yaml.Features = make([]string, 0, 5)
	for _, f := range p.templates_def.Features.Common {
		if v, err := Evaluate(f, p.Model()); err != nil {
			return fmt.Errorf("Unable to evaluate '%s'. %s", f, err)
		} else {
			if v == "" {
				log.Printf("INFO! No feature defined with '%s'.", f)
				continue
			}
			f = v
		}
		if f != "" {
			p.yaml.Features = append(p.yaml.Features, f)
		}
	}

	if deploy_features, ok := p.templates_def.Features.Deploy[p.yaml.Deploy.Deployment.To]; ok {
		for _, f := range deploy_features {
			if f != "" {
				p.yaml.Features = append(p.yaml.Features, f)
			}
		}
	}

	p.CleanModel()

	// TODO: Load additionnal features from maintainer source path or file. This will permit adding more features and let the plugin manage generated path from update task.

	// Load all sources
	p.sources = make(map[string]TmplSource)
	p.templates = make(map[string]TmplSource)

	for file, f := range p.templates_def.Sources.Common {
		if f.Template == "" {
			p.sources[file] = f
			log.Printf("%#v", p.sources[file])
		} else {
			p.templates[file] = f
		}
	}

	if deploy_sources, ok := p.templates_def.Sources.Deploy[p.yaml.Deploy.Deployment.To]; ok {
		for file, f := range deploy_sources {
			if f.Template == "" {
				p.sources[file] = f
			} else {
				p.templates[file] = f
			}
		}
	}
	log.Printf("Files: \n%#v\n\n%#v\n", p.sources, p.templates)

	return nil
}

func (ts *TmplSource) Generate(tmpl_data interface{}, template_dir, dest_path, dest_name string) (updated bool, _ error) {
	src := path.Join(template_dir, ts.Template)
	dest := path.Join(dest_path, dest_name)
	parent := path.Dir(dest)

	if parent != "." {
		if _, err := os.Stat(parent); err != nil {
			os.MkdirAll(parent, 0755)
			updated = true
		}
	}

	var data string
	if b, err := ioutil.ReadFile(src); err != nil {
		return false, fmt.Errorf("Load issue. %s", err)
	} else {
		data = strings.Replace(string(b), "}}\\\n", "}}", -1)
	}

	t, err := template.New(src).Funcs(template.FuncMap{}).Parse(data)
	if err != nil {
		return false, fmt.Errorf("Template issue. %s", err)
	}

	orig_md5, _ := md5sum(dest)
	final_md5_file := md5.New()

	if out, err := os.Create(dest); err != nil {
		return false, fmt.Errorf("Unable to create %s. %s.", dest, err)
	} else {
		multi_write_file := io.MultiWriter(out, final_md5_file)
		if err := t.Execute(multi_write_file, tmpl_data); err != nil {
			return false, fmt.Errorf("Unable to interpret %s. %s.", dest, err)
		}
		out.Close()
	}
	final_md5 := final_md5_file.Sum(nil)
	if orig_md5 != nil {
		updated = updated || !bytes.Equal(orig_md5, final_md5)
	} else {
		updated = true
	}

	if u, err := set_rights(dest, ts.Chmod); err != nil {
		return false, fmt.Errorf("%s", err)
	} else {
		updated = updated || u
	}
	return
}
