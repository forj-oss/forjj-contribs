package main

import (
    "github.hpe.com/christophe-larsonneur/goforjj"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "path"
)

type JenkinsPlugin struct {
    yaml YamlJenkins
    source_path string
    template_dir string
    template_file string
    templates_def YamlTemplates // See templates.go
    data TemplateData
    sources map[string]TmplSource
    templates map[string]TmplSource
}

type TemplateData struct {
     Features []string
}

type YamlJenkins struct {
    Source SourceStruct
    Deploy DeployStruct
    Features []string
}


const jenkins_file = "forjj-jenkins.yaml"


func new_plugin(src string) (p *JenkinsPlugin) {
    p = new(JenkinsPlugin)

    p.source_path = src
    return
}

func (p *JenkinsPlugin) initialize_from(r *CreateReq, ret *goforjj.PluginData) (status bool) {
    p.yaml.Source = r.SourceStruct
    p.yaml.Deploy = r.DeployStruct
    return true
}

func (p *JenkinsPlugin) load_from(ret *goforjj.PluginData) (status bool) {
    return true
}

func (p *JenkinsPlugin) update_from(r *UpdateReq, ret *goforjj.PluginData)  (status bool) {
    return true
}

func (p *JenkinsPlugin)save_yaml(ret *goforjj.PluginData) (status bool) {
    file := path.Join(p.yaml.Source.ForjjInstanceName, jenkins_file)

    d, err := yaml.Marshal(&p.yaml)
    if  err != nil {
        ret.Errorf("Unable to encode forjj-jenkins configuration data in yaml. %s", err)
        return
    }

    if err := ioutil.WriteFile(file, d, 0644) ; err != nil {
        ret.Errorf("Unable to save '%s'. %s", file, err)
        return
    }
    ret.AddFile(file)
    return true
}

func (p *JenkinsPlugin)load_yaml(ret *goforjj.PluginData) (status bool) {

    return true
}
