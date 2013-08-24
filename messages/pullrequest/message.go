package pullrequest

import (
	"github.com/mhrivnak/ghreleaseguard/messages"
)

type Links struct {
	Self struct {
		Href string
	}
}

type PullRequest struct {
	Base struct {
		Ref string
	}
	Links  Links `json:"_links"`
	Number uint
}

type Message struct {
	PullRequest PullRequest `json:"pull_request"`
	Repository  struct {
		Owner struct {
			Login string
		}
		Name string
	}
}

// Version examines the Message's "Ref" attribute and returns a version string,
// if found.
func (message *Message) Version() (string, bool) {
	return messages.Version(message.PullRequest.Base.Ref)
}
