package pullrequest

import (
	"github.com/mhrivnak/ghreleaseguard/messages"
)

type Commit struct {
	Id string
}

type Href struct {
	Href string
}

type Links struct {
	Self Href
}

type Owner struct {
	Login string
}

type Repository struct {
	Owner Owner
	Name  string
}

type Base struct {
	Ref string
}

type PullRequest struct {
	Base   Base
	Links  Links `json:"_links"`
	Number uint
}

type Message struct {
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository
}

// Version examines the Message's "Ref" attribute and returns a version string,
// if found.
func (message *Message) Version() (string, bool) {
	return messages.Version(message.PullRequest.Base.Ref)
}
