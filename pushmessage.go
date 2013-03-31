package main

import (
	"regexp"
)

var releaseExp *regexp.Regexp

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

func initRegexp() {
	var err error
	releaseExp, err = regexp.Compile(`.*(\d+\.\d+)$`)
	if err != nil {
		panic(err)
	}
}
