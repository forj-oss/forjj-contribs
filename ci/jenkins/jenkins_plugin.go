package main

import (
    "github.com/forj-oss/goforjj"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "path"
    "log"
    "fmt"
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

type DeployApp struct {
    DeployStruct `yaml:",inline"`
    Command string // Command to use to execute a Deploy
}

type ForjjStruct struct {
    InstanceName string
    OrganizationName string
}

// Used for the jenkins yaml source and generate template data.
type YamlJenkins struct {
    Forjj ForjjStruct
    // Settings SettingsStruct
    Deploy DeployApp
    Features []string
    Dockerfile DockerfileStruct
    JenkinsImage FinalImageStruct
}

/*type SettingsStruct struct {
}*/


const jenkins_file = "forjj-jenkins.yaml"


func new_plugin(src string) (p *JenkinsPlugin) {
    p = new(JenkinsPlugin)

    p.source_path = src
    return
}

// At create time: create jenkins source from req
func (p *JenkinsPlugin) initialize_from(r *CreateReq, ret *goforjj.PluginData) (status bool) {
    p.yaml.Forjj.InstanceName = r.Args.ForjjInstanceName
    p.yaml.Forjj.OrganizationName = r.Args.ForjjOrganization
    p.yaml.Deploy.DeployStruct = r.Args.DeployStruct
    // Forjj predefined settings (instance/organization) are set at create time only.
    // I do not recommend to update them, manually by hand in the `forjj-jenkins.yaml`.
    // Updating the instance name could be possible but not for now.
    // As well Moving an instance to another orgnization could be possible, but I do not see a real use case.
    // So, they are fixed and saved at create time. Update/maintain won't never update them later.
    if err := p.DefineDeployCommand() ; err != nil {
        ret.Errorf("Unable to define deployement command. %s", err)
        return
    }

    p.yaml.Dockerfile = r.Args.DockerfileStruct
    p.yaml.JenkinsImage.SetFrom(&r.Args.FinalImageStruct, r.Args.ForjjOrganization)
    return true
}

func (p *JenkinsPlugin) DefineDeployCommand() error{
    if err := p.LoadTemplatesDef() ; err != nil {
        return fmt.Errorf("%s", err)
    }

    if v, ok := p.templates_def.Run[p.yaml.Deploy.DeployTo] ; !ok {
        list := make([]string,0,len(p.templates_def.Run))
        for element, _ := range p.templates_def.Run {
            list = append(list, element)
        }
        return fmt.Errorf("'%s' deploy type is unknown (templates.yaml). Valid are %s", p.yaml.Deploy.DeployTo, list)
    } else {
        p.yaml.Deploy.Command = v
    }
    return nil
}

// TODO: Detect if the commands was manually updated to avoid updating it if end user did it alone.

// At update time: Update jenkins source from req or forjj-jenkins.yaml input.
func (p *JenkinsPlugin) update_from(r *UpdateReq, ret *goforjj.PluginData)  (status bool) {
    // ForjjStruct NOT UPDATABLE
    p.yaml.Deploy.SetFrom(&r.Args.DeployStruct)
    if err := p.DefineDeployCommand() ; err != nil {
        ret.Errorf("Unable to update the deployement command. %s", err)
        return
    }
    p.yaml.Dockerfile.SetFrom(&r.Args.DockerfileStruct)
    p.yaml.JenkinsImage.SetFrom(&r.Args.FinalImageStruct, r.Args.ForjjOrganization)// Org used only if no set anymore.
    return true
}

func (p *JenkinsPlugin)save_yaml(ret *goforjj.PluginData) (status bool) {
    file := path.Join(p.source_path, jenkins_file)

    d, err := yaml.Marshal(&p.yaml)
    if  err != nil {
        ret.Errorf("Unable to encode forjj-jenkins configuration data in yaml. %s", err)
        return
    }

    if err := ioutil.WriteFile(file, d, 0644) ; err != nil {
        ret.Errorf("Unable to save '%s'. %s", file, err)
        return
    }
    // Be careful to not introduce the local mount which in containers can be totally different (due to docker -v)
    ret.AddFile(path.Join(p.yaml.Forjj.InstanceName, jenkins_file))
    ret.StatusAdd("'%s' instance saved (%s).", p.yaml.Forjj.InstanceName, path.Join(p.yaml.Forjj.InstanceName, jenkins_file))
    log.Printf("'%s' instance saved.", file)
    return true
}

func (p *JenkinsPlugin)load_yaml(ret *goforjj.PluginData) (status bool) {
    file := path.Join(p.source_path, jenkins_file)

    log.Printf("Loading '%s'...", file)
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
