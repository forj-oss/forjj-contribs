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
    Deployments map[string]DeployStruct
	// Those 2 different parameters are defined at create time and can be updated with change.
	// There are default deployment task and name. This can be changed at maintain time
	// to reflect the maintain deployment task to execute.
	DeployTo string  // Default Deployment set at create time.
	Command string   // Default Command used
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
	instance := r.Forj.ForjjInstanceName
    p.yaml.Forjj.InstanceName = instance
    p.yaml.Forjj.OrganizationName = r.Forj.ForjjOrganization

	if _, found := r.Objects.App[instance] ; !found {
		ret.Errorf("Request format issue. Unable to find the jenkins instance '%s'", instance)
		return
	}
	jenkins_instance := r.Objects.App[instance]

	// Initialize deployment data and set default values
	if p.yaml.Deploy.Deployments == nil {
		p.yaml.Deploy.Deployments = make(map[string]DeployStruct)
	}
	if len(r.Objects.Deployment) == 0 {
		// Set default deployment with docker.
		p.yaml.Deploy.Deployments["docker"] = DeployStruct{
			Name: "docker",
			ServiceAddr: "localhost",
			ServicePort: "8080",
			DeployTo: "docker",
		}
		ret.StatusAdd("Added default docker deployment.")
	} else {
		// Set Deployments definition
		for name, deploy_data := range r.Objects.Deployment {
			deployment := DeployStruct{}
			deployment.SetFrom(&deploy_data.Add.AddDeployStruct)
			p.yaml.Deploy.Deployments[name] = deployment
		}
	}

    deploy_to := jenkins_instance.Add.DeployTo
	if deploy_to == "" {
		deploy_to = "docker"
		ret.StatusAdd("default deployment to 'docker'.")
	}

	if _, found := p.yaml.Deploy.Deployments[deploy_to] ; !found {
		ret.Errorf("Unable to find deployment '%s'. You must define it.", deploy_to)
	}

	// Default deployment set
	p.yaml.Deploy.DeployTo = deploy_to

    // Forjj predefined settings (instance/organization) are set at create time only.
    // I do not recommend to update them, manually by hand in the `forjj-jenkins.yaml`.
    // Updating the instance name could be possible but not for now.
    // As well Moving an instance to another organization could be possible, but I do not see a real use case.
    // So, they are fixed and saved at create time. Update/maintain won't never update them later.
    if err := p.DefineDeployCommand() ; err != nil {
        ret.Errorf("Unable to define the default deployement command. %s", err)
        return
    }

	// Initialize Dockerfile data and set default values
	p.yaml.Dockerfile.SetFrom(&jenkins_instance.Add.AddDockerfileStruct)

	// Initialize Jenkins Image data and set default values
	p.yaml.JenkinsImage.SetFrom(&jenkins_instance.Add.AddFinalImageStruct, r.Forj.ForjjOrganization)

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
	instance := r.Forj.ForjjInstanceName
	instance_data := r.Objects.App[instance].Change
	if d, found := r.Objects.Deployment[instance_data.DeployTo] ; !found {
		ret.Errorf("Unable to find deployment '%s'", instance_data.DeployTo)
		return
	} else {
		deploy := DeployStruct{}
		if _, found := r.Objects.Deployment[instance_data.DeployTo] ; !found {
			deploy.SetFrom(&d.Add.AddDeployStruct)
			ret.StatusAdd("Deployment '%s' added.", instance_data.DeployTo)
		} else {
			deploy = p.yaml.Deploy.Deployments[instance_data.DeployTo]
			deploy.UpdateFrom(&d.Change.ChangeDeployStruct)
			ret.StatusAdd("Deployment '%s' updated.", instance_data.DeployTo)
		}
		p.yaml.Deploy.Deployments[instance_data.DeployTo] = deploy
	}

    if err := p.DefineDeployCommand() ; err != nil {
        ret.Errorf("Unable to update the deployement command. %s", err)
        return
    }

    p.yaml.Dockerfile.UpdateFrom(&instance_data.ChangeDockerfileStruct)
    p.yaml.JenkinsImage.UpdateFrom(&instance_data.ChangeFinalImageStruct, r.Forj.ForjjOrganization)// Org used only if no set anymore.
	status = true
    return
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
