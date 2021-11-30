package version

import (
	"regexp"
	"github.com/Masterminds/semver/v3"
)

var rcRegExp = regexp.MustCompile(`rc[1-9]\d*`)

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

func IsDev(version string) (res bool, err error) {
	v, err := parse(version)
	if err != nil {
		return
	}
	res = "dev" == v.Prerelease()
	return
}

func IsRc(version string) (res bool, err error) {
	v, err := parse(version)
	if err != nil {
		return
	}
	res = rcRegExp.MatchString(v.Prerelease())
	return
}

//func NextDev(version string) (res string, err error) {
//	v, err := parse(version)
//	if err != nil {
//		return
//	}
//	res = ""
//	return
//}
//
//func NextRc(version string) (res string, err error) {
//	v, err := parse(version)
//	if err != nil {
//		return
//	}
//	res = ""
//	return
//}

