package main

import (
    "github.hpe.com/christophe-larsonneur/goforjj"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "path"
    "log"
)

type JenkinsPlugin struct {
    yaml YamlJenkins
    source_path string
    template_dir string
    template_file string
    templates_def YamlTemplates // See templates.go
    sources map[string]TmplSource
    templates map[string]TmplSource
}

type DockerStruct struct {
    Name string
    Version string
    Repository string
    Maintainer string
}

type DeployApp struct {
    DeployStruct `yaml:",inline"`
    DeployCommand string // Command to use to execute a Deploy
}

type ForjjStruct struct {
    InstanceName string
    OrganizationName string
}

// Used for the jenkins yaml source and generate template data.
type YamlJenkins struct {
    Forjj ForjjStruct
    // Settings SettingsStruct
    Docker DockerStruct
    Deploy DeployApp
    Features []string
}

/*type SettingsStruct struct {
}*/


const jenkins_file = "forjj-jenkins.yaml"


func new_plugin(src string) (p *JenkinsPlugin) {
    p = new(JenkinsPlugin)

    p.source_path = src
    return
}

// Update jenkins source from input sources
func (p *JenkinsPlugin) initialize_from(r *CreateReq, ret *goforjj.PluginData) (status bool) {
    p.yaml.Docker.SetFrom(&r.Args.SourceStruct, r.Args.ForjjOrganization)
    p.yaml.Deploy.DeployStruct = r.Args.DeployStruct
    p.yaml.Forjj.InstanceName = r.Args.ForjjInstanceName
    p.yaml.Forjj.OrganizationName = r.Args.ForjjOrganization
    // Forjj predefined settings (instance/organization) are set at create time only.
    // I do not recommend to update them, manually by hand in the `forjj-jenkins.yaml`.
    // Updating the instance name could be possible but not for now.
    // As well Moving an instance to another orgnization could be possible, but I do not see a real use case.
    // So, they are fixed and saved at create time. Update/maintain won't never update them later.
    return true
}

func (p *JenkinsPlugin) load_from(ret *goforjj.PluginData) (status bool) {
    return true
}

func (p *JenkinsPlugin) update_from(r *UpdateReq, ret *goforjj.PluginData)  (status bool) {
    p.yaml.Deploy.SetFrom(&r.Args.DeployStruct)
    p.yaml.Docker.SetFrom(&r.Args.SourceStruct, "") // No orga name to provide for updates...
    return true
}

func (p *JenkinsPlugin)save_yaml(ret *goforjj.PluginData) (status bool) {
    file := path.Join(p.yaml.Forjj.InstanceName, jenkins_file)

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
    log.Printf("'%s' instance saved.")
    return true
}

func (p *JenkinsPlugin)load_yaml(instance string, ret *goforjj.PluginData) (status bool) {
    file := path.Join(instance, jenkins_file)

    if d, err := ioutil.ReadFile(file) ; err != nil {
        ret.Errorf("Unable to read '%s'. %s", file, err)
        return
    } else {
        if  err = yaml.Unmarshal(d, &p.yaml) ; err != nil {
            ret.Errorf("Unable to decode forjj-jenkins configuration data from yaml. %s", err)
            return
        }
    }
    log.Printf("'%s' instance loaded.", file)
    return true
}
