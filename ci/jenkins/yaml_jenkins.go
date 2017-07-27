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
