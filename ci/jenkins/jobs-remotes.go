package main

type JobRemote interface {
	GetUpstream() string
}
