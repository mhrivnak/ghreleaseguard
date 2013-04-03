package handlers

import (
	"encoding/json"
	"github.com/mhrivnak/ghreleaseguard/config"
	"github.com/mhrivnak/ghreleaseguard/messages/pullrequest"
	"io/ioutil"
	"log"
	"net/http"
)

func inspectPullRequest(raw []byte) {
	var message pullrequest.Message
	err := json.Unmarshal(raw, &message)
	if err != nil {
		log.Println("push.inspect: ", err)
		return
	}
	log.Println(message)

	// if we can't find a version, quit early
	versionName, found := message.Version()
	if !found {
		log.Println("version not found")
		return
	}

	// if we can't find a forbidden commit, quit early
	forbiddenCommit, found := config.GetForbiddenCommit(
		message.Repository.Owner.Login,
		message.Repository.Name,
		versionName)
	if !found {
		return
	}
	log.Println("Found forbidden commit: ", forbiddenCommit)
}

func PullRequestHandler(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("PullRequestHandler: ", err)
		return
	}
	go inspectPullRequest(raw)
}
