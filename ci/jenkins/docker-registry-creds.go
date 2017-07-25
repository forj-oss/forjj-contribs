package main

import (
	"encoding/base64"
	"fmt"
	"github.com/forj-oss/forjj-modules/trace"
	"log"
	"regexp"
	"strings"
)

type DockerAuths struct {
	Auths map[string]DockerRegistryCredsInfo
}

type DockerRegistryCredsInfo struct {
	user     string
	password string
}

func (a *DockerAuths) authenticate(server string) error {
	if _, found := a.Auths[server]; !found {
		return fmt.Errorf("Unable to authenticate to docker registry '%s'. Server not found.", server)
	}
	auth := a.Auths[server]
	gotrace.SetDebug()
	if cmdlog, err := run_cmd("sudo", nil, "docker", "login", "-u", auth.user, "-p", auth.password, server); err != nil {
		log.Printf("Unable to authenticate. %s. docker output: %s", err, cmdlog)
		return err
	} else {
		log.Printf("%s", cmdlog)
	}
	return nil
}

func NewDockerAuths(auths string) (a *DockerAuths) {
	auths_s := strings.Split(auths, ",")

	auth_reg, _ := regexp.Compile(`([^:]+):(\w+=*)`)

	a = new(DockerAuths)

	a.Auths = make(map[string]DockerRegistryCredsInfo)

	for _, v := range auths_s {
		if substr := auth_reg.FindStringSubmatch(v); substr != nil {
			var user_pwd []string
			if v, err := base64.StdEncoding.DecodeString(substr[2]); err != nil {
				log.Printf("Unable to decode base64 '%s'", substr[2])
				return
			} else {
				user_pwd = strings.Split(strings.TrimSpace(string(v)), ":")
			}
			auth := DockerRegistryCredsInfo{user_pwd[0], user_pwd[1]}
			a.Auths[substr[1]] = auth
			log.Printf("'%s' registry server credential added.", substr[1])
		}
	}
	return
}
