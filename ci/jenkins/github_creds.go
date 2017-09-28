package main

import "fmt"

type UserPasswordCreds struct {
	Name     string `yaml:"name,omitempty"`
	password string
}

func (t *UserPasswordCreds) SetFrom(d *GithubUserStruct) (status bool) {
	status = SetIfSet(&t.Name, d.Username)
	return
}

func (t *UserPasswordCreds) UpdateFrom(d *GithubUserStruct) (status bool) {
	status = SetOrClean(&t.Name, d.Username)
	return
}

func (t *UserPasswordCreds) setPassword(password string) (_ error) {
	if t.Name == "" && password == "" {
		return
	}
	if t.Name == "" {
		return fmt.Errorf("You set the github user password, but the github user name is missing. " +
			"Please, update your Forjfile.")
	}
	if password == "" {
		return fmt.Errorf("Password for '%s' is missing. Please set the github password and retry.", t.Name)
	}
	t.password = password
	return
}

func (t *UserPasswordCreds) GetPassword() string {
	return t.password
}
