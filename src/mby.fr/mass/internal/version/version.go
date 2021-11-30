package version

import (

	"github.com/Masterminds/semver/v3"
)

func parse(version string) (v *semver.Version, err error) {
	//v, err = semver.StrictNewVersion(version)
	v, err = semver.NewVersion(version)
	return
}

func NextPatch(version string) (res string, err error) {
	v, err := parse(version)
	if err != nil {
		return
	}
	bumped := v.IncPatch()
	res = bumped.String()
	return
}

func NextMinor(version string) (res string, err error) {
	v, err := parse(version)
	if err != nil {
		return
	}
	bumped := v.IncMinor()
	res = bumped.String()
	return
}

func NextMajor(version string) (res string, err error) {
	v, err := parse(version)
	if err != nil {
		return
	}
	bumped := v.IncMajor()
	res = bumped.String()
	return
}

