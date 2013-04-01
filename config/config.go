package config

import (
	"code.google.com/p/gcfg"
	"log"
)

// This data structure gets traversed to determine if a given branch has a
// forbidden commit defined
var Owners map[string]Owner

type Owner struct {
	Repos map[string]Repo
}

type Repo struct {
	Versions map[string]string
}

// directly represents a parsed config entry
type release struct {
	Commit  string
	Owner   string
	Repo    string
	Version string
}

// internal storage of the config data as loaded from disk
type config struct {
	Release map[string]*release
}

// true iff the config has been loaded successfully
var configLoaded bool

// LoadConfig loads the config from disk, causes it to be parsed, and then
// makes available the data structure Owners. If there is an error parsing the
// file, but the config has previously been loaded, this will leave the existing
// config data in place and log a message about the failure.
func LoadConfig() {
	var values config
	err := gcfg.ReadFileInto(&values, "/etc/ghreleaseguard.conf")
	if err != nil {
		if !configLoaded {
			panic(err)
		} else {
			// if we already have a valid config, just keep using it
			log.Println("Failed to load config: ", err)
		}
	}
	Owners = parseConfig(values)
	configLoaded = true
}

// parseConfig accepts the parsed JSON data as structs and builds a data
// structure suitable for the "Owners" value.
func parseConfig(values config) map[string]Owner {
	ret := make(map[string]Owner)
	for _, release := range values.Release {
		owner, ok := ret[release.Owner]
		if !ok {
			owner = Owner{make(map[string]Repo)}
			ret[release.Owner] = owner
		}
		repo, ok := owner.Repos[release.Repo]
		if !ok {
			repo = Repo{make(map[string]string)}
			owner.Repos[release.Repo] = repo
		}
		repo.Versions[release.Version] = release.Commit
	}
	return ret
}

// GetForbiddenCommit traverses the "Owners" data structure based on data about
// a branch and returns a forbidden commit ID if one is found.
func GetForbiddenCommit(ownerName, repoName, version string) (string, bool) {
	owner, ok := Owners[ownerName]
	if !ok {
		return "", false
	}
	repo, ok := owner.Repos[repoName]
	if !ok {
		return "", false
	}
	commit, ok := repo.Versions[version]
	return commit, ok
}
