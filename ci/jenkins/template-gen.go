package main

func (t *DockerStruct)SetFrom(d *SourceStruct, org string) {
    SetIfSet(&t.Name, d.DockerImage)
    SetIfSet(&t.Version, d.DockerImageVersion)
    SetIfSet(&t.Repository, d.DockerRepoimage)
    SetOnceIfSet(&t.Repository, org)
    SetIfSet(&t.Maintainer, d.Maintainer)
}

func (t *DeployStruct)SetFrom(d *DeployStruct) {
    SetIfSet(&t.DeployTo, d.DeployTo)
    SetIfSet(&t.ServiceAddr, d.ServiceAddr)
    SetIfSet(&t.ServicePort, d.ServicePort)
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
