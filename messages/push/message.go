package push

import (
	"github.com/mhrivnak/ghreleaseguard/messages"
)

// below types are used for receiving JSON data from GitHub

type User struct {
	Email string
	Name  string
}

type Message struct {
	Commits []struct {
		Id string
	}
	Pusher     User
	Ref        string
	Repository struct {
		Name  string
		Owner User
		Url   string
	}
}

// Version examines the Message's "Ref" attribute and returns a version string,
// if found.
func (message *Message) Version() (string, bool) {
	return messages.Version(message.Ref)
}
