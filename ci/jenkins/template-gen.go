package main

func (t *DeployStruct)SetFrom(d *AddDeployStruct) {
    SetIfSet(&t.DeployTo, d.DeployTo)
    SetIfSet(&t.ServiceAddr, d.ServiceAddr)
    SetIfSet(&t.ServicePort, d.ServicePort)
}

func (t *DockerfileStruct)SetFrom(d *AddDockerfileStruct) {
    SetIfSet(&t.FromImage, d.FromImage)
    SetIfSet(&t.FromImageVersion, d.FromImageVersion)
    SetIfSet(&t.Maintainer, d.Maintainer)
}

func (t *FinalImageStruct)SetFrom(d *FinalImageStruct, org string) {
    SetIfSet(&t.Name, d.Name)
    SetIfSet(&t.FinalDockerImageVersion, d.FinalDockerImageVersion)
    SetIfSet(&t.RegistryServer, d.RegistryServer)

    SetIfSet(&t.FinalDockerRepoName, d.FinalDockerRepoName)
    SetOnceIfSet(&t.FinalDockerRepoName, org)
}

// Set the value if the source is set
func SetIfSet(s *string, source string) {
    if source == "" {
        return
    }
    *s = source
}

// Set the value originally empty from source if set.
func SetOnceIfSet(s *string, source string) {
    if *s != "" || source == "" {
        return
    }
    *s = source
}
