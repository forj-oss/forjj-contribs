package main

import (
    "github.hpe.com/christophe-larsonneur/goforjj"
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
func InstantiateAll(mount string, ret *goforjj.PluginData) (status bool) {
    elements, err := ioutil.ReadDir(".")
    if err != nil {
        ret.Errorf("Issue to read Current directory. %s", elements)
        return false
    }

    for _, instance := range elements {
        src := path.Join(mount, instance.Name())
        if _, err := os.Stat(path.Join(src, jenkins_file)) ; err == nil {
            p := new_plugin(src)
            ret.StatusAdd("Maintaining '%s'", instance.Name())
            if err := os.Chdir(src) ; err != nil {
                ret.Errorf("Unable to enter in '%s'. %s", src, err)
                return false
            }
            if ! p.InstantiateInstance(ret) {
                return false
            }
        } else {
            log.Printf("'%s' is not a forjj plugin source code model. No '%s' found. ignored.", src, jenkins_file)
        }
    }
    return true
}

func (p *JenkinsPlugin)InstantiateInstance(ret *goforjj.PluginData) (status bool) {
    if ! p.load_yaml(ret) {
        return
    }

    // start a command as described by the source code.
    if p.yaml.Deploy.Command == "" {
        ret.Errorf("Unable to instantiate to %s. Deploy Command is empty.", p.yaml.Deploy.DeployTo)
        return
    }
    ret.StatusAdd("Running '%s'", p.yaml.Deploy.Command)
    s, err := run_cmd("/bin/sh", "-c", p.yaml.Deploy.Command)
    ret.StatusAdd(string(s))
    if err != nil {
        cur_dir, _ := os.Getwd()
        ret.Errorf("%s (pwd: %s)", err, cur_dir)
    }

    return true
}
