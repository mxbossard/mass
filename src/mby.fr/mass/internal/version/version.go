package version

import (
	"log"

	"github.com/Masterminds/semver/v3"
)

func BumpPatch(version string) *semver.Version {
	v, err := semver.NewVersion(version)
	if err != nil {
		log.Fatal(err)
	}
	bumped := v.IncPatch()
	return &bumped
}
