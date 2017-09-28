package main

// Used for the jenkins yaml source and generate template data.
type YamlJenkins struct {
	Forjj ForjjStruct
	// Settings SettingsStruct
	Deploy       DeployApp
	Features     []string
	Dockerfile   DockerfileStruct
	JenkinsImage FinalImageStruct
	Projects     *Projects
	admin_pwd    string
	GithubUser   UserPasswordCreds
}

func (y *YamlJenkins) ProjectsHasSource(name string) (_ bool) {
	if y == nil || y.Projects == nil {
		return
	}
	for _, project := range y.Projects.List {
		if project.SourceType == name {
			return true
		}
	}
	return
}

func (y YamlJenkins) GetAdminPwd() string {
	return y.admin_pwd
}

func (y *YamlJenkins) SetAdminPwd(pwd string) {
	y.admin_pwd = pwd
}
