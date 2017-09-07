package main

import (
	"github.com/forj-oss/goforjj"
	"log"
	"os"
	"path"
)

// Return ok if the jenkins instance exist
func (r *MaintainReq) check_source_existence(ret *goforjj.PluginData) (status bool) {
	log.Print("Checking Jenkins source code path existence.")

	src_path := path.Join(r.Forj.ForjjSourceMount, r.Forj.ForjjInstanceName)
	if _, err := os.Stat(path.Join(src_path, jenkins_file)); err != nil {
		log.Printf(ret.Errorf("Unable to maintain instance name '%s' without source code.\n"+
			"Use update to update it, commit, push and retry. %s.", src_path, err))
		return
	}

	ret.StatusAdd("environment checked.")
	status = true
	return
}

// TODO: Need to define where to deploy (dev/itg/pro/local/other) - Is it still needed?

// Instantiate Instance given by the request.
func (r *MaintainReq) Instantiate(req *MaintainReq, ret *goforjj.PluginData) (_ bool) {
	instance := r.Forj.ForjjInstanceName
	mount := r.Forj.ForjjSourceMount
	auths := NewDockerAuths(r.Objects.App[instance].RegistryAuth)

	src := path.Join(mount, instance)
	if _, err := os.Stat(path.Join(src, jenkins_file)); err == nil {
		p := new_plugin(src)
		if ! p.GetMaintainData(instance, req, ret) {
			return false
		}
		ret.StatusAdd("Maintaining '%s'", instance)
		if err := os.Chdir(src); err != nil {
			ret.Errorf("Unable to enter in '%s'. %s", src, err)
			return
		}
		if !p.InstantiateInstance(instance, auths, ret) {
			return false
		}
	} else {
		log.Printf("'%s' is not a forjj plugin source code model. No '%s' found. ignored.", src, jenkins_file)
	}
	return true
}

func (p *JenkinsPlugin) InstantiateInstance(instance string, auths *DockerAuths, ret *goforjj.PluginData) (status bool) {
	if !p.load_yaml(ret) {
		return
	}

	run, found := p.templates_def.Run[p.yaml.Deploy.Deployment.To]
	if ! found {
		ret.Errorf("Deployment '%s' command not found.")
		return
	}

	// start a command as described by the source code.
	if run.RunCommand == "" {
		log.Printf(ret.Errorf("Unable to instantiate to %s. Deploy Command is empty.", p.yaml.Deploy.Deployment.To))
		return
	}

	for server := range auths.Auths {
		if err := auths.authenticate(server); err != nil {
			log.Printf(ret.Errorf("Unable to instantiate. %s", err))
			return
		}
	}

	log.Printf(ret.StatusAdd("Running '%s'", run.RunCommand))

	var env []string
	if v := os.Getenv("DOOD_SRC"); v != "" {
		env = append(os.Environ(), "SRC="+path.Join(v, instance))
	}

	model := p.Model()
	for key, env_to_set := range run.Env {
		if env_to_set.If != "" {
			// check if If evaluation return something or not. if not, the environment key won't be created.
			if v, err := Evaluate(env_to_set.If, model) ; err != nil {
				ret.Errorf("Deployment '%s'. Error in evaluating '%s'. %s", p.yaml.Deploy.Deployment.To, key, err)
			} else {
				if v == "" {
					continue
				}
			}
		}
		if v, err := Evaluate(env_to_set.Value, model) ; err != nil {
			ret.Errorf("Deployment '%s'. Error in evaluating '%s'. %s", p.yaml.Deploy.Deployment.To, key, err)
		} else {
			env = append(os.Environ(), key + "=" + v)
			log.Printf("Env added : '%s' = '%s'", key, v)
		}
	}

	s, err := run_cmd("/bin/sh", env, "-c", run.RunCommand)
	log.Printf(ret.StatusAdd(string(s)))
	if err != nil {
		cur_dir, _ := os.Getwd()
		log.Printf(ret.Errorf("%s (pwd: %s)", err, cur_dir))
	}

	return true
}
