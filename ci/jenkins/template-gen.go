package main

func (t *DeployStruct) SetFrom(d *DeployStruct) (status bool) {
	if t == nil {
		return false
	}
	status = SetIfSet(&t.ServiceAddr, d.ServiceAddr)
	status = SetOnceIfSet(&t.To, d.To) || status
	return SetIfSet(&t.ServicePort, d.ServicePort) || status
}

func (t *DeployStruct) UpdateFrom(d *DeployStruct) (status bool) {
	status = SetOrClean(&t.ServiceAddr, d.ServiceAddr)
	status = SetOnceIfSet(&t.To, d.To) || status
	return SetOrClean(&t.ServicePort, d.ServicePort) || status
}

func (t *YamlSSLStruct) UpdateFrom(d *SslStruct) (status bool) {
	status = SetOrClean(&t.Certificate, d.Certificate) || status
	return
}

func (t *YamlSSLStruct) SetFrom(d *SslStruct) (status bool) {
	status = SetIfSet(&t.Certificate, d.Certificate) || status
	return
}

func (t *DockerfileStruct) SetFrom(d *DockerfileStruct) (status bool) {
	status = SetIfSet(&t.FromImage, d.FromImage)
	status = SetIfSet(&t.FromImageVersion, d.FromImageVersion) || status
	return SetIfSet(&t.Maintainer, d.Maintainer) || status
}

func (t *DockerfileStruct) UpdateFrom(d *DockerfileStruct) (status bool) {
	status = SetOrClean(&t.FromImage, d.FromImage)
	status = SetOrClean(&t.FromImageVersion, d.FromImageVersion) || status
	return SetOrClean(&t.Maintainer, d.Maintainer) || status
}

func (t *FinalImageStruct) SetFrom(d *FinalImageStruct, org string) (status bool) {
	status = SetIfSet(&t.Name, d.Name)
	status = SetIfSet(&t.Version, d.Version) || status
	status = SetIfSet(&t.RegistryServer, d.RegistryServer) || status

	status = SetIfSet(&t.RegistryRepoName, d.RegistryRepoName) || status
	return SetOnceIfSet(&t.RegistryRepoName, org) || status
}

func (t *FinalImageStruct) UpdateFrom(d *FinalImageStruct, org string) (status bool) {
	status = SetOrClean(&t.Name, d.Name)
	status = SetOrClean(&t.Version, d.Version) || status
	status = SetOrClean(&t.RegistryServer, d.RegistryServer) || status

	status = SetOrClean(&t.RegistryRepoName, d.RegistryRepoName) || status
	return SetOnceIfSet(&t.RegistryRepoName, org) || status
}

func (t YamlSSLStruct) GetKey() string {
	return t.key
}

func (t *YamlSSLStruct) SetKey(key string) bool {
	return SetIfSet(&t.key, key)
}

// SetIfSet Set the value if the source is set
// return true if updated.
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

// SetOrClean simply copy the value
// return true if updated.
func SetOrClean(s *string, source string) (_ bool) {
	if *s != source {
		*s = source
		return true
	}
	return
}

// SetOnceIfSet Set the value originally empty from source if set.
// return true if updated.
func SetOnceIfSet(s *string, source string) (_ bool) {
	if *s != "" || source == "" {
		return
	}
	if *s != source {
		*s = source
		return true
	}
	return
}
