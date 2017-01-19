package main

import (
    "github.com/forj-oss/goforjj"
    "log"
    "os"
    "path"
    "io/ioutil"
)

// Return ok if the jenkins instance exist
func (r *MaintainReq) check_source_existence(ret *goforjj.PluginData) (status bool) {
    log.Printf("Checking Jenkins source code path existence.")

    if _, err := os.Stat(r.Args.ForjjSourceMount) ; err != nil {
        ret.Errorf("Unable to maintain jenkins instances. '%s' is inexistent or innacessible. %s", r.Args.ForjjSourceMount, err)
        return
    }

    ret.StatusAdd("environment checked.")
    status = true
    return
}

// Looping on all instances
// Need to define where to deploy (dev/itg/pro/local/other)
func (r *MaintainReq)InstantiateAll(ret *goforjj.PluginData) (status bool) {
    elements, err := ioutil.ReadDir(".")
    if err != nil {
        ret.Errorf("Issue to read Current directory. %s", elements)
        return false
    }

    mount := r.Args.ForjjSourceMount
    auths := NewDockerAuths(r.Args.RegistryAuth)

    for _, instance := range elements {
        src := path.Join(mount, instance.Name())
        if _, err := os.Stat(path.Join(src, jenkins_file)) ; err == nil {
            p := new_plugin(src)
            ret.StatusAdd("Maintaining '%s'", instance.Name())
            if err := os.Chdir(src) ; err != nil {
                ret.Errorf("Unable to enter in '%s'. %s", src, err)
                return
            }
            if ! p.InstantiateInstance(instance.Name(), auths, ret) {
                return false
            }
        } else {
            log.Printf("'%s' is not a forjj plugin source code model. No '%s' found. ignored.", src, jenkins_file)
        }
    }
    return true
}

func (p *JenkinsPlugin)InstantiateInstance(instance string, auths *DockerAuths, ret *goforjj.PluginData) (status bool) {
    if ! p.load_yaml(ret) {
        return
    }

    // start a command as described by the source code.
    if p.yaml.Deploy.Command == "" {
        ret.Errorf("Unable to instantiate to %s. Deploy Command is empty.", p.yaml.Deploy.DeployTo)
        return
    }

    if err := auths.write_docker_config() ; err != nil {
        ret.Errorf("Unable to instantiate. Unable to write the docker registry credential file. %s", err)
        return
    }
    defer auths.remove_config()

    ret.StatusAdd("Running '%s'", p.yaml.Deploy.Command)

    var env []string
    if v := os.Getenv("DOOD_SRC") ; v != "" {
        env = append(os.Environ(), "SRC=" + path.Join(v, instance))
    }

    s, err := run_cmd("/bin/sh", env, "-c", p.yaml.Deploy.Command)
    ret.StatusAdd(string(s))
    if err != nil {
        cur_dir, _ := os.Getwd()
        ret.Errorf("%s (pwd: %s)", err, cur_dir)
    }

    return true
}
