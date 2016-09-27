package main

func (t *DeployStruct)SetFrom(d *DeployStruct) {
    SetIfSet(&t.DeployTo, d.DeployTo)
    SetIfSet(&t.ServiceAddr, d.ServiceAddr)
    SetIfSet(&t.ServicePort, d.ServicePort)
}

func (t *DockerfileStruct)SetFrom(d *DockerfileStruct) {
    SetIfSet(&t.BaseDockerImage, d.BaseDockerImage)
    SetIfSet(&t.BaseDockerImageVersion, d.BaseDockerImageVersion)
    SetIfSet(&t.Maintainer, d.Maintainer)
}

func (t *FinalImageStruct)SetFrom(d *FinalImageStruct, org string) {
    SetIfSet(&t.FinalDockerImage, d.FinalDockerImage)
    SetIfSet(&t.FinalDockerImageVersion, d.FinalDockerImageVersion)
    SetIfSet(&t.FinalDockerRegistryServer, d.FinalDockerRegistryServer)

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
