package main

func (t *DeployStruct)SetFrom(d *DeployStruct) (status bool) {
	if t == nil {
		return false
	}
    status = SetIfSet(&t.ServiceAddr, d.ServiceAddr)
    return SetIfSet(&t.ServicePort, d.ServicePort) || status
}

func (t *DeployStruct)UpdateFrom(d *DeployStruct) (status bool) {
    status = SetIfSet(&t.ServiceAddr, d.ServiceAddr)
	return SetIfSet(&t.ServicePort, d.ServicePort) || status
}

func (t *DockerfileStruct)SetFrom(d *DockerfileStruct) (status bool) {
    status = SetIfSet(&t.FromImage, d.FromImage)
	status = SetIfSet(&t.FromImageVersion, d.FromImageVersion) || status
	return SetIfSet(&t.Maintainer, d.Maintainer) || status
}

func (t *DockerfileStruct)UpdateFrom(d *DockerfileStruct) (status bool) {
	status = SetIfSet(&t.FromImage, d.FromImage)
	status = SetIfSet(&t.FromImageVersion, d.FromImageVersion) || status
	return SetIfSet(&t.Maintainer, d.Maintainer) || status
}

func (t *FinalImageStruct)SetFrom(d *FinalImageStruct, org string) (status bool) {
    status = SetIfSet(&t.Name, d.Name)
	status = SetIfSet(&t.Version, d.Version) || status
	status = SetIfSet(&t.RegistryServer, d.RegistryServer) || status

	status = SetIfSet(&t.RegistryRepoName, d.RegistryRepoName) || status
	return SetOnceIfSet(&t.RegistryRepoName, org) || status
}

func (t *FinalImageStruct)UpdateFrom(d *FinalImageStruct, org string) (status bool) {
	status = SetIfSet(&t.Name, d.Name)
	status = SetIfSet(&t.Version, d.Version) || status
	status = SetIfSet(&t.RegistryServer, d.RegistryServer) || status

	status = SetIfSet(&t.RegistryRepoName, d.RegistryRepoName) || status
	return SetOnceIfSet(&t.RegistryRepoName, org) ||status
}

// SetIfSet Set the value if the source is set
func SetIfSet(s *string, source string) (_ bool) {
    if source == "" {
        return
    }
	if *s != source {
		*s = source
		return true
	}
    return
}

// SetOnceIfSet Set the value originally empty from source if set.
func SetOnceIfSet(s *string, source string) (_ bool){
    if *s != "" || source == "" {
        return
    }
	if *s != source {
		*s = source
		return true
	}
	return
}
