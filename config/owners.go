package config

// This data structure gets traversed to identify a forbidden commit based on an
// owner, repo, and version
var Owners map[string]Owner

type Owner struct {
	Repos map[string]Repo
}

type Repo struct {
	Versions map[string]string
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
