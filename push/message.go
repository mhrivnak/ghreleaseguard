package push

import (
	"regexp"
)

// below types are used for receiving JSON data from GitHub

type Commit struct {
	Id string
}

type User struct {
	Email string
	Name  string
}

type Repository struct {
	Name  string
	Owner User
}

type Message struct {
	Commits    []Commit
	Pusher     User
	Ref        string
	Repository Repository
}

var versionExp = regexp.MustCompile(`.*-(\d+\.\d+)$`)

// Version examines the Message's "Ref" attribute and returns a version string,
// if found.
func (message *Message) Version() (string, bool) {
	result := versionExp.FindStringSubmatch(message.Ref)
	if len(result) == 2 {
		return result[1], true
	}
	return "", false
}
