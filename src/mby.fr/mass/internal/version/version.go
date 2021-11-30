package version

import (
	"log"

	"github.com/Masterminds/semver/v3"
)

func BumpPatch(version string) string {
	v, err := semver.NewVersion(version)
	if err != nil {
		log.Printf("Unable to parse version: [%v] ! %v", version, err)
		return ""
	} else {
		log.Printf("Parsed version: [%v].", v)
	}
	bumped := v.IncPatch()
	log.Printf("Bumped version: [%v].", bumped)
	return bumped.String()
}
