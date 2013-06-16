package handlers

import (
	"encoding/json"
	"github.com/mhrivnak/ghreleaseguard/config"
	"github.com/mhrivnak/ghreleaseguard/messages/pullrequest"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
)

func inspectPullRequest(raw []byte) {
	var message pullrequest.Message
	err := json.Unmarshal(raw, &message)
	if err != nil {
		log.Println("push.inspect: ", err)
		return
	}

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

	commits, err := getPRCommits(message.PullRequest.Links.Self.Href)
	if err != nil {
		return
	}

	for _, commit := range commits {
		if commit.Sha == forbiddenCommit {
			log.Println("MATCH! forbidden commit is in the PR")
		}
	}
}

func getPRCommits(href string) ([]Commit, error) {
	// API call to get commits in this PR
	commitURL, err := url.Parse(href)
	if err != nil {
		log.Println("error parsing PR URL: ", err)
		return nil, err
	}
	commitURL.Path = path.Join(commitURL.Path, "commits")
	response, err := http.Get(commitURL.String())
	if err != nil {
		log.Println("error getting commits: ", err)
		return nil, err
	}
	defer response.Body.Close()

	rawCommits, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("error reading commits response: ", err)
		return nil, err
	}
	var commits []Commit
	err = json.Unmarshal(rawCommits, &commits)
	if err != nil {
		log.Println("error parsing commit json: ", err)
		return nil, err
	}
	return commits, nil
}

type Commit struct {
	Sha string
}

func PullRequestHandler(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("PullRequestHandler: ", err)
		return
	}
	go inspectPullRequest(raw)
}
