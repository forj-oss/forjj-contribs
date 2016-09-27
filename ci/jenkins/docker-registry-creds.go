package main

import (
    "encoding/json"
    "os"
    "log"
    "strings"
    "regexp"
    "fmt"
)

type DockerAuths struct {
    Auths map[string]DockerRegistryCredsInfo
}

type DockerRegistryCredsInfo struct {
    Auth string
    Email string `json:",omitempty"`
}

const protected_config = "/tmp/docker_config.json"

func (a *DockerAuths)write_docker_config() error {
    if out, err := os.Create(protected_config) ; err != nil {
        return fmt.Errorf("Unable to create %s. %s.", protected_config, err)
    } else {
        defer out.Close()
        if err := json.NewEncoder(out).Encode(a) ; err != nil {
            return fmt.Errorf("Unable to generate the docker registry credential file '%s'. %s.", protected_config, err)
        }
    }
    os.Chmod(protected_config, 0600)
    log.Printf("'%s' file generated.", protected_config)

    if cmdlog, err := run_cmd("sudo", "/bin/docker-config-update.sh") ; err != nil {
        log.Printf("Unable to update docker config file. %s. Script output: %s", err, cmdlog)
        return err
    } else {
        log.Printf("%s", cmdlog)
    }
    return nil
}

func (a *DockerAuths)remove_config() {
    //os.Remove(protected_config)
    log.Printf("'%s' file Removed.", protected_config)
}

func NewDockerAuths(auths string) (a *DockerAuths) {
    auths_s := strings.Split(auths, ",")

    auth_reg, _ := regexp.Compile(`([^:]+):(\w+)(:(.*@.*))?`)

    a = new(DockerAuths)

    a.Auths = make(map[string]DockerRegistryCredsInfo)

    for _, v := range auths_s {
        if substr := auth_reg.FindStringSubmatch(v) ; substr != nil {
            auth := DockerRegistryCredsInfo{substr[2], substr[4]}
            a.Auths[substr[1]] = auth
            log.Printf("'%s' registry server credential added.", substr[1])
        }
    }
    return
}
