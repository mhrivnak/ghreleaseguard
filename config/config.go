package config

import (
	"code.google.com/p/gcfg"
	"log"
	"os"
)

// directly represents a parsed config entry
type release struct {
	Commit  string
	Owner   string
	Repo    string
	Version string
}

// internal storage of the config data as loaded from disk
type config struct {
	Server  Server
	Release map[string]*release
}

// true iff the config has been loaded successfully
var configLoaded bool

// LoadConfig loads the config from disk, causes it to be parsed, and then
// makes available the ServerConfig and Owners data structures. If there is an
// error parsing the file, but the config has previously been loaded, this will
// leave the existing config data in place and log a message about the failure.
func LoadConfig() {
	var values config
	err := gcfg.ReadFileInto(&values, getConfigPath())
	if err != nil {
		if !configLoaded {
			panic(err)
		} else {
			// if we already have a valid config, just keep using it
			log.Println("Failed to load config: ", err)
		}
	}
	ServerConfig = values.Server
	Owners = values.Owners()
	configLoaded = true
}

// looks in an environment variable for a path to the config file, and if not
// found, returns the default.
func getConfigPath() string {
	path := os.Getenv("GHRGCONFIGPATH")
	if path == "" {
		return "/etc/ghreleaseguard.conf"
	} else {
		return path
	}
}

// parseConfig operates on the parsed config data as structs and builds a data
// structure suitable for the "Owners" value.
func (values *config) Owners() map[string]Owner {
	ret := make(map[string]Owner)
	for _, release := range values.Release {
		// owner may or may not already be in the return map
		owner, ok := ret[release.Owner]
		if !ok {
			owner = Owner{make(map[string]Repo)}
			ret[release.Owner] = owner
		}
		// repo may or may not already be in the owner map
		repo, ok := owner.Repos[release.Repo]
		if !ok {
			repo = Repo{make(map[string]string)}
			owner.Repos[release.Repo] = repo
		}
		repo.Versions[release.Version] = release.Commit
	}
	return ret
}
