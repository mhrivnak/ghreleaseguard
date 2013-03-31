package main

import (
	"regexp"
)

var releaseExp = regexp.MustCompile(`.*-(\d+\.\d+)$`)

type Commit struct {
	Id string
}

type Pusher struct {
	Email string
	Name  string
}

type PushMessage struct {
	Commits []Commit
	Pusher  Pusher
	Ref     string
}

func (message *PushMessage) Release() string {
	result := releaseExp.FindStringSubmatch(message.Ref)
	if len(result) == 2 {
		return result[1]
	}
	return ""
}
