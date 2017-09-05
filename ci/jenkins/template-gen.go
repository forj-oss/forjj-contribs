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
	status = SetIfSet(&t.ServiceAddr, d.ServiceAddr)
	status = SetOnceIfSet(&t.To, d.To) || status
	return SetIfSet(&t.ServicePort, d.ServicePort) || status
}

func (t *YamlSSLStruct) UpdateFrom(d *SslStruct) bool {
	return t.SetFrom(d)
}

func (t *YamlSSLStruct) SetFrom(d *SslStruct) (status bool) {
	status = SetIfSet(&t.CaCertificate, d.CaCertificate)
	status = SetIfSet(&t.Certificate, d.Certificate) || status
	return
}

func (t *DockerfileStruct) SetFrom(d *DockerfileStruct) (status bool) {
	status = SetIfSet(&t.FromImage, d.FromImage)
	status = SetIfSet(&t.FromImageVersion, d.FromImageVersion) || status
	return SetIfSet(&t.Maintainer, d.Maintainer) || status
}

func (t *DockerfileStruct) UpdateFrom(d *DockerfileStruct) (status bool) {
	status = SetIfSet(&t.FromImage, d.FromImage)
	status = SetIfSet(&t.FromImageVersion, d.FromImageVersion) || status
	return SetIfSet(&t.Maintainer, d.Maintainer) || status
}

func (t *FinalImageStruct) SetFrom(d *FinalImageStruct, org string) (status bool) {
	status = SetIfSet(&t.Name, d.Name)
	status = SetIfSet(&t.Version, d.Version) || status
	status = SetIfSet(&t.RegistryServer, d.RegistryServer) || status

	status = SetIfSet(&t.RegistryRepoName, d.RegistryRepoName) || status
	return SetOnceIfSet(&t.RegistryRepoName, org) || status
}

func (t *FinalImageStruct) UpdateFrom(d *FinalImageStruct, org string) (status bool) {
	status = SetIfSet(&t.Name, d.Name)
	status = SetIfSet(&t.Version, d.Version) || status
	status = SetIfSet(&t.RegistryServer, d.RegistryServer) || status

	status = SetIfSet(&t.RegistryRepoName, d.RegistryRepoName) || status
	return SetOnceIfSet(&t.RegistryRepoName, org) || status
}

func (t *DeployApp) SetFrom(source *DeployApp) (status bool) {
	if t == nil {
		return
	}

	status = SetIfSet(&t.Command, source.Command)
	status = SetIfSet(&t.Ssl.Certificate, source.Ssl.Certificate) || status
	status = SetIfSet(&t.Ssl.CaCertificate, source.Ssl.CaCertificate) || status
	return
}

func (t *YamlSSLStruct) GetKey() string {
	return t.key
}

func (t *YamlSSLStruct) SetKey(key string) bool {
	return SetIfSet(&t.key, key)
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
