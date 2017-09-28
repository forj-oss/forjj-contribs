package main

import "fmt"

type UserPasswordCreds struct {
	UserName string `yaml:"user_name,omitempty"`
	password string
}

func (t *UserPasswordCreds) SetFrom(d *GithubCredStruct) (status bool) {
	status = SetIfSet(&t.UserName, d.Username)
	return
}

func (t *UserPasswordCreds) UpdateFrom(d *GithubCredStruct) (status bool) {
	status = SetOrClean(&t.UserName, d.Username)
	return
}

func (t *UserPasswordCreds) setPassword(password string) (_ error) {
	if t.UserName == "" && password == "" {
		return
	}
	if t.UserName == "" {
		return fmt.Errorf("You set the github user password, but the github user name is missing. " +
			"Please, update your Forjfile.")
	}
	if password == "" {
		return fmt.Errorf("Password for '%s' is missing. Please set the github password and retry.", t.UserName)
	}
	t.password = password
	return
}

func (t *UserPasswordCreds) GetPassword() string {
	return t.password
}
