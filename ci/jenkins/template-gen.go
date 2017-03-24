package main

func (t *DeployStruct)SetFrom(d *DeployStruct) bool {
	if t == nil {
		return false
	}
    SetIfSet(&t.DeployTo, d.DeployTo)
    SetIfSet(&t.ServiceAddr, d.ServiceAddr)
    SetIfSet(&t.ServicePort, d.ServicePort)
	return true
}

func (t *DeployStruct)UpdateFrom(d *DeployStruct) {
    SetIfSet(&t.DeployTo, d.DeployTo)
    SetIfSet(&t.ServiceAddr, d.ServiceAddr)
    SetIfSet(&t.ServicePort, d.ServicePort)
}

func (t *DockerfileStruct)SetFrom(d *DockerfileStruct) {
    SetIfSet(&t.FromImage, d.FromImage)
    SetIfSet(&t.FromImageVersion, d.FromImageVersion)
    SetIfSet(&t.Maintainer, d.Maintainer)
}

func (t *DockerfileStruct)UpdateFrom(d *DockerfileStruct) {
	SetIfSet(&t.FromImage, d.FromImage)
	SetIfSet(&t.FromImageVersion, d.FromImageVersion)
	SetIfSet(&t.Maintainer, d.Maintainer)
}

func (t *FinalImageStruct)SetFrom(d *FinalImageStruct, org string) {
    SetIfSet(&t.Name, d.Name)
    SetIfSet(&t.Version, d.Version)
    SetIfSet(&t.RegistryServer, d.RegistryServer)

    SetIfSet(&t.RegistryRepoName, d.RegistryRepoName)
    SetOnceIfSet(&t.RegistryRepoName, org)
}

func (t *FinalImageStruct)UpdateFrom(d *FinalImageStruct, org string) {
	SetIfSet(&t.Name, d.Name)
	SetIfSet(&t.Version, d.Version)
	SetIfSet(&t.RegistryServer, d.RegistryServer)

	SetIfSet(&t.RegistryRepoName, d.RegistryRepoName)
	SetOnceIfSet(&t.RegistryRepoName, org)
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
