package main


import (
    "os"
    "gopkg.in/yaml.v2"
    "fmt"
    "io/ioutil"
    "log"
)

const template_file = "templates.yaml"

// Contains functions to manage source code from templates

type YamlTemplates struct {
    Features TmplFeatures
    Sources TmplSources
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

func (p *JenkinsPlugin)DefineSources() error {
    // load all features
    p.data.Features = make([]string, 0, 5)
    for _, f := range p.templates_def.Features.Common {
        p.data.Features = append(p.data.Features, f)
    }

    if deploy_features, ok:= p.templates_def.Features.Deploy[p.yaml.Deploy.DeployTo] ; ok {
        for _, f := range deploy_features {
            p.data.Features = append(p.data.Features, f)
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

    if deploy_sources, ok:= p.templates_def.Sources.Deploy[p.yaml.Deploy.DeployTo] ; ok {
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
