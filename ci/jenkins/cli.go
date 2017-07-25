// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"github.com/forj-oss/goforjj"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"log"
)

type JenkinsApp struct {
	App    *kingpin.Application
	params Params
	socket string
	Yaml   goforjj.YamlPlugin
}

type Params struct {
	socket_file  *string
	socket_path  *string
	template_dir *string
	daemon       *bool // Currently not used - Lot of concerns with daemonize in go... Stay in foreground
}

func (a *JenkinsApp) init() {
	a.load_plugin_def()

	a.App = kingpin.New("jenkins", "CI jenkins plugin for FORJJ.")
	version := "0.1"
	if version != "" {
		a.App.Version(version)
	}

	// true to create the Infra
	daemon := a.App.Command("service", "jenkins REST API service")
	daemon.Command("start", "start jenkins REST API service")
	a.params.socket_file = daemon.Flag("socket-file", "Socket file to use").Default(a.Yaml.Runtime.Service.Socket).String()
	a.params.socket_path = daemon.Flag("socket-path", "Socket file path to use").Default("/tmp/forjj-socks").String()
	a.params.daemon = daemon.Flag("daemon", "Start process in background like a daemon").Short('d').Bool()
	a.params.template_dir = daemon.Flag("templates", "Path to templates files.").Default("templates").String()

}

func (a *JenkinsApp) load_plugin_def() {
	yaml.Unmarshal([]byte(YamlDesc), &a.Yaml)
	if a.Yaml.Runtime.Service.Socket == "" {
		a.Yaml.Runtime.Service.Socket = "jenkins.sock"
		log.Printf("Set default socket file: %s", a.Yaml.Runtime.Service.Socket)
	}
}
