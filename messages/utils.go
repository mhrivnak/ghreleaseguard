package messages

import (
	"regexp"
)

var versionExp = regexp.MustCompile(`.*-(\d+\.\d+)$`)

// Version examines the Message's "Ref" attribute and returns a version string,
// if found.
func Version(ref string) (string, bool) {
	result := versionExp.FindStringSubmatch(ref)
	if len(result) == 2 {
		return result[1], true
	}
	return "", false
}
