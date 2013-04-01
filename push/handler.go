package push

import (
	"encoding/json"
	"github.com/mhrivnak/ghreleaseguard/config"
	"io/ioutil"
	"log"
	"net/http"
)

// inspect parses the raw report sent by github, determines if a forbidden
// commit is present in the push, and takes action if so.
func inspect(raw []byte) {
	var message Message
	err := json.Unmarshal(raw, &message)
	if err != nil {
		log.Println("push.inspect: ", err)
		return
	}

	// if we can't find a version, quit early
	versionName, found := message.Version()
	if !found {
		return
	}

	// if we can't find a forbidden commit, quit early
	forbiddenCommit, found := config.GetForbiddenCommit(
		message.Repository.Owner.Name,
		message.Repository.Name,
		versionName)
	if !found {
		return
	}
	log.Println("Found forbidden commit: ", forbiddenCommit)

	// search this push for the forbidden commit
	for _, commit := range message.Commits {
		if commit.Id == forbiddenCommit {
			// Take action here! probably send an email
			log.Println("MATCH! forbidden commit is in the push")
		}
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("push.Handler: ", err)
		return
	}
	go inspect(raw)
}
