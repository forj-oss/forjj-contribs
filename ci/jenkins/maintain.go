package main

import (
    "github.hpe.com/christophe-larsonneur/goforjj"
    "log"
    "os"
    "path"
    "io/ioutil"
)

// Return ok if the jenkins instance exist
func (r *MaintainReq) check_source_existence(ret *goforjj.PluginData) (p *JenkinsPlugin, status bool) {
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
func (p *JenkinsPlugin)InstantiateAll(ret *goforjj.PluginData) (status bool) {
    elements, err := ioutil.ReadDir(".")
    if err != nil {
        ret.Errorf("Issue to read Current directory. %s", elements)
        return false
    }

    for _, instance := range elements {
        if _, err := os.Stat(path.Join(instance.Name(), jenkins_file)) ; err == nil {
            if ! p.InstantiateInstance(instance.Name(), ret) {
                return false
            }
        }
    }
    return true
}

func (p *JenkinsPlugin)InstantiateInstance(instance string, ret *goforjj.PluginData) (status bool) {
    ret.StatusAdd("Maintaining '%s'", instance)
    if ! p.load_yaml(instance, ret) {
        return
    }

    // start a command as described by the source code.
    if p.yaml.Deploy.DeployCommand == "" {
        ret.Errorf("Unable to instantiate to %s. Deploy Command is empty.", p.yaml.Deploy.DeployTo)
        return
    }
    ret.StatusAdd("Running '%s'", p.yaml.Deploy.DeployCommand)
    s, err := run_cmd("/bin/sh", "-c", p.yaml.Deploy.DeployCommand)
    ret.StatusAdd(string(s))
    if err != nil {
        ret.Errorf("%s", err)
    }

    return true
}
