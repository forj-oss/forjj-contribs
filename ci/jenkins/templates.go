package main


import (
    "os"
    "gopkg.in/yaml.v2"
    "fmt"
    "io/ioutil"
    "log"
    "strings"
    "path"
    "text/template"
)

const template_file = "templates.yaml"

// Contains functions to manage source code from templates

type YamlTemplates struct {
    Defaults DefaultsStruct
    Features TmplFeatures
    Sources TmplSources
    Run map[string]string `yaml:"run_deploy"`
    Variants map[string]string
}

type DefaultsStruct struct {
	Dockerfile DockerfileStruct
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
    Chmod os.FileMode
    Template string
    Source string
}

//Load templates definition file from template dir.
func (p *JenkinsPlugin)LoadTemplatesDef() error {
    if d, err := ioutil.ReadFile(p.template_file) ; err != nil {
        return fmt.Errorf("Unable to load '%s'. %s.", p.template_file, err)
    } else {
        if err := yaml.Unmarshal(d, &p.templates_def) ; err != nil {
            return fmt.Errorf("Unable to load yaml file format '%s'. %s.", p.template_file, err)
        }
    }
    return nil
}

// Load list of files to copy and files to generate
func (p *JenkinsPlugin)DefineSources() error {
    // load all features
    p.yaml.Features = make([]string, 0, 5)
    for _, f := range p.templates_def.Features.Common {
        p.yaml.Features = append(p.yaml.Features, f)
    }

    if deploy_features, ok:= p.templates_def.Features.Deploy[p.yaml.Deploy.Deployment.To] ; ok {
        for _, f := range deploy_features {
            p.yaml.Features = append(p.yaml.Features, f)
        }
    }

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

    if deploy_sources, ok:= p.templates_def.Sources.Deploy[p.yaml.Deploy.Deployment.To] ; ok {
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

func (ts *TmplSource)Generate(tmpl_data interface{}, template_dir, dest_path, dest_name string) error {
    src := path.Join(template_dir, ts.Template)
    dest := path.Join(dest_path, dest_name)
    parent := path.Dir(dest)

    if  parent != "." {
        if _, err := os.Stat(parent) ; err != nil {
            os.MkdirAll(parent, 0755)
        }
    }

    var data string
    if b, err := ioutil.ReadFile(src) ; err != nil {
        return fmt.Errorf("Load issue. %s", err)
    } else {
        data = strings.Replace(string(b), "}}\\\n", "}}", -1)
    }

    t, err := template.New(src).Funcs(template.FuncMap{}).Parse(data)
    if err != nil {
        return fmt.Errorf("Template issue. %s", err)
    }

    if out, err := os.Create(dest) ; err != nil {
        return fmt.Errorf("Unable to create %s. %s.", dest, err)
    } else {
        if err := t.Execute(out, tmpl_data) ; err != nil {
            return fmt.Errorf("Unable to interpret %s. %s.", dest, err)
        }
        out.Close()
    }

    if err := set_rights(dest, ts.Chmod) ; err != nil {
        return fmt.Errorf("%s", err)
    }
    return nil
}
