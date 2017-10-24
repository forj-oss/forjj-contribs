package main

import (
	"fmt"
	"github.com/forj-oss/goforjj"
	"net/url"
	"regexp"
)

type ProjectInfo struct {
	ForjCommonStruct
	infra_remote string
}

func (pi *ProjectInfo) set_project_info(forj ForjCommonStruct) {
	pi.ForjCommonStruct = forj
}

func (pi *ProjectInfo) set_infra_remote(infra_remote string) {
	pi.infra_remote = infra_remote
}

func (pi *ProjectInfo) set_projects_to(projects map[string]ProjectsInstanceStruct, r *JenkinsPlugin,
	ret *goforjj.PluginData, status *bool, InfraName string) (_ error) {
	if pi.ForjjInfraUpstream == "" {
		ret.StatusAdd("Unable to add a new project without a remote GIT repository. Jenkins JobDSL requirement. " +
			"To enable this feature, add a remote GIT to your infra --infra-upstream or define the JobDSL Repository to clone.")
		IsUpdated(status)
		return
	}

	ssh_format, _ := regexp.Compile(`^(https?://)(\w[\w.-]+)((/(\w[\w.-]*)/(\w[\w.-]*))(/\w[\w.-/]*)?)$`)
	job_path := ""
	default_jobdsl := false
	if rs := ssh_format.FindStringSubmatch(pi.infra_remote); rs != nil {
		if rs[5] == pi.ForjjOrganization && rs[6] == pi.ForjjInfra {
			job_path = "jobs-dsl"
			default_jobdsl = true
		} else {
			pi.infra_remote = rs[1] + rs[2] + rs[4]
			job_path = rs[7]
		}
	}

	if v, err := url.Parse(pi.infra_remote); err != nil {
		ret.Errorf("Infra remote URL issue. %s", err)
		return err
	} else {
		if v.Scheme == "" {
			err = fmt.Errorf("Invalid Remote repository Url '%s'. A scheme must exist.", pi.infra_remote)
			ret.Errorf("%s", err)
			return err
		}
	}
	// Initialize JobDSL structure
	r.yaml.Projects = NewProjects(pi.ForjjInstanceName, pi.infra_remote, job_path, default_jobdsl)

	// Retrieve list of Repository (projects) to manage
	for name, prj := range projects {
		switch prj.RemoteType {
		case "github":
			r.yaml.Projects.AddGithub(name, &prj.GithubStruct, (name == InfraName))
		case "git":
			r.yaml.Projects.AddGit(name, &prj.GitStruct, (name == InfraName))
		}
	}
	IsUpdated(status)
	return
}
