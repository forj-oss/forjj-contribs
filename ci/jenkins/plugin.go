package main

//go:generate go get github.hpe.com/christophe-larsonneur/goforjj gopkg.in/yaml.v2
//go:generate go build -o $GOPATH/bin/forjj-genapp github.hpe.com/christophe-larsonneur/goforjj/genapp
//go:generate forjj-genapp jenkins.yaml

import (
    "github.hpe.com/christophe-larsonneur/goforjj"
)

type JenkinsPlugin struct {
    yaml YamlJenkins
    source_path string
    template_dir string
    template_file string
    templates_def YamlTemplates // See templates.go
    data TemplateData
    sources map[string]*TmplSource
    templates map[string]*TmplSource
}

type TemplateData struct {
     Features []string
}

type YamlJenkins struct {
    Source SourceStruct
    Deploy DeployStruct
    Features []string
}

func (r *GroupReq)new_plugin(src string) (p *JenkinsPlugin) {
    p = &JenkinsPlugin{}

    p.source_path = src
    return
}

func (p *JenkinsPlugin) initialize_from(r *CreateReq, ret *goforjj.PluginData) (status bool) {
    p.yaml.Source = r.Groups.Source
    p.yaml.Deploy = r.Groups.Deploy
    return true
}

func (p *JenkinsPlugin) load_from(ret *goforjj.PluginData) (status bool) {
    return true
}

func (p *JenkinsPlugin) update_from(r *UpdateReq, ret *goforjj.PluginData)  (status bool) {
    return true
}

func (r *JenkinsPlugin)save_yaml(ret *goforjj.PluginData) (status bool) {

    return true
}

func (r *JenkinsPlugin)load_yaml(ret *goforjj.PluginData) (status bool) {

    return true
}
