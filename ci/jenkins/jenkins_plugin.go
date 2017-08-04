package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/forj-oss/goforjj"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
)

type JenkinsPluginModel struct {
	Source YamlJenkins
}

var JP_Model *JenkinsPluginModel

type JenkinsPlugin struct {
	yaml          YamlJenkins
	source_path   string
	template_dir  string
	template_file string
	templates_def YamlTemplates // See templates.go
	sources       map[string]TmplSource
	templates     map[string]TmplSource
}

type DeployApp struct {
	Deployment DeployStruct
	// Those 2 different parameters are defined at create time and can be updated with change.
	// There are default deployment task and name. This can be changed at maintain time
	// to reflect the maintain deployment task to execute.
	Command string // Default Command used
}

type ForjjStruct struct {
	InstanceName     string
	OrganizationName string
	InfraUpstream    string
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
	instance := r.Forj.ForjjInstanceName
	p.yaml.Forjj.InstanceName = instance
	p.yaml.Forjj.OrganizationName = r.Forj.ForjjOrganization
	p.yaml.Forjj.InfraUpstream = r.Forj.ForjjInfraUpstream

	if _, found := r.Objects.App[instance]; !found {
		ret.Errorf("Request format issue. Unable to find the jenkins instance '%s'", instance)
		return
	}
	jenkins_instance := r.Objects.App[instance]

	p.yaml.Deploy.Deployment.SetFrom(&jenkins_instance.DeployStruct)
	// Initialize deployment data and set default values
	if p.yaml.Deploy.Deployment.To == "" {
		p.yaml.Deploy.Deployment.To = "docker"
		ret.StatusAdd("Default to 'docker' Deployment.")
	}
	if p.yaml.Deploy.Deployment.ServiceAddr == "" {
		p.yaml.Deploy.Deployment.ServiceAddr = "localhost"
		ret.StatusAdd("Default to 'localhost' deployment service name.")
	}
	if p.yaml.Deploy.Deployment.ServicePort == "" {
		p.yaml.Deploy.Deployment.ServicePort = "8080"
		ret.StatusAdd("Default to '8080' deployment service port.")

	}

	// Forjj predefined settings (instance/organization) are set at create time only.
	// I do not recommend to update them, manually by hand in the `forjj-jenkins.yaml`.
	// Updating the instance name could be possible but not for now.
	// As well Moving an instance to another organization could be possible, but I do not see a real use case.
	// So, they are fixed and saved at create time. Update/maintain won't never update them later.
	if err := p.DefineDeployCommand(); err != nil {
		ret.Errorf("Unable to define the default deployment command. %s", err)
		return
	}

	// Initialize Dockerfile data and set default values
	log.Printf("CreateReq : %#v\n", r)
	p.yaml.Dockerfile.SetFrom(&jenkins_instance.DockerfileStruct)
	log.Printf("p.yaml.Dockerfile : %#v\n", p.yaml.Dockerfile)

	// Initialize Jenkins Image data and set default values
	p.yaml.JenkinsImage.SetFrom(&jenkins_instance.FinalImageStruct, r.Forj.ForjjOrganization)

	if !p.add_projects(r, ret) {
		return
	}

	status = true
	return
}

func (p *JenkinsPlugin) DefineDeployCommand() error {
	if err := p.LoadTemplatesDef(); err != nil {
		return fmt.Errorf("%s", err)
	}

	if v, ok := p.templates_def.Run[p.yaml.Deploy.Deployment.To]; !ok {
		list := make([]string, 0, len(p.templates_def.Run))
		for element, _ := range p.templates_def.Run {
			list = append(list, element)
		}
		return fmt.Errorf("'%s' deploy type is unknown (templates.yaml). Valid are %s", p.yaml.Deploy.Deployment.To, list)
	} else {
		p.yaml.Deploy.Command = v
	}
	return nil
}

// TODO: Detect if the commands was manually updated to avoid updating it if end user did it alone.

// At update time: Update jenkins source from req or forjj-jenkins.yaml input.
func (p *JenkinsPlugin) update_from(r *UpdateReq, ret *goforjj.PluginData) (status bool) {
	instance := r.Forj.ForjjInstanceName
	instance_data := r.Objects.App[instance]
	deploy := DeployStruct{}
	deploy = p.yaml.Deploy.Deployment
	if status = deploy.UpdateFrom(&instance_data.DeployStruct); status {
		ret.StatusAdd("Deployment to '%s' updated.", instance_data.To)
	}
	p.yaml.Deploy.Deployment = deploy

	if err := p.DefineDeployCommand(); err != nil {
		ret.Errorf("Unable to update the deployement command. %s", err)
		return
	}

	if p.yaml.Dockerfile.UpdateFrom(&instance_data.DockerfileStruct) {
		ret.StatusAdd("Dockerfile updated.")
		status = true
	}
	// Org used only if no set anymore.
	if p.yaml.JenkinsImage.UpdateFrom(&instance_data.FinalImageStruct, r.Forj.ForjjOrganization) {
		ret.StatusAdd("Jenkins master docker image data updated.")
		status = true
	}
	return
}

func (p *JenkinsPlugin) save_yaml(ret *goforjj.PluginData) (status bool) {
	file := path.Join(p.source_path, jenkins_file)

	orig_md5, _ := md5sum(file)
	d, err := yaml.Marshal(&p.yaml)
	if err != nil {
		ret.Errorf("Unable to encode forjj-jenkins configuration data in yaml. %s", err)
		return
	}
	final_md5 := md5.New().Sum(d)

	if bytes.Equal(orig_md5, final_md5) {
		return false
	}

	if err := ioutil.WriteFile(file, d, 0644); err != nil {
		ret.Errorf("Unable to save '%s'. %s", file, err)
		return
	}
	// Be careful to not introduce the local mount which in containers can be totally different (due to docker -v)
	ret.AddFile(path.Join(p.yaml.Forjj.InstanceName, jenkins_file))
	ret.StatusAdd("'%s' instance saved (%s).", p.yaml.Forjj.InstanceName, path.Join(p.yaml.Forjj.InstanceName, jenkins_file))
	log.Printf("'%s' instance saved.", file)
	return true
}

func (p *JenkinsPlugin) load_yaml(ret *goforjj.PluginData) (status bool) {
	file := path.Join(p.source_path, jenkins_file)

	log.Printf("Loading '%s'...", file)
	if d, err := ioutil.ReadFile(file); err != nil {
		ret.Errorf("Unable to read '%s'. %s", file, err)
		return
	} else {
		if err = yaml.Unmarshal(d, &p.yaml); err != nil {
			ret.Errorf("Unable to decode forjj-jenkins configuration data from yaml. %s", err)
			return
		}
	}
	log.Printf("'%s' instance loaded.", file)
	return true
}
