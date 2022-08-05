package version

import (
	"fmt"
	"regexp"
	"strconv"
	"github.com/Masterminds/semver/v3"
)

var rcRegExp = regexp.MustCompile(`rc(\d+)`)

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

func NextDev(version string) (res string, err error) {
	v, err := parse(version)
	if err != nil {
		return
	}

	if ok, err := IsDev(version); ok {
		return version, err
	}
	if ok, _ := IsRc(version); ok {
		nextVer, _ := NextPatch(version)
		return NextDev(nextVer)
	}

	bumped := v.IncPatch()
	bumped, err = bumped.SetPrerelease("dev")
	res = bumped.String()
	return
}

func NextRc(version string) (res string, err error) {
	v, err := parse(version)
	if err != nil {
		return
	}

	var nth uint32 = 1
	bumped := *v
	if ok, _ := IsRc(version); ok {
		// Increment rc qualifier
		rcQualifier := v.Prerelease()
		submatches := rcRegExp.FindStringSubmatch(rcQualifier)
		n, err := strconv.ParseUint(submatches[1], 10, 32)
		if err != nil {
			return "", err
		}
		nth = uint32(n)
		nth ++
	} else {
		// Increment patch
		bumped = bumped.IncPatch()
	}

	bumped, err = bumped.SetPrerelease("rc" + fmt.Sprint(nth))
	res = bumped.String()
	return
}


func Dev(version string) (res string, err error) {
	v, err := parse(version)
	if err != nil {
		return
	}

	bumped, err := v.SetPrerelease("dev")
	res = bumped.String()
	return
}

func Release(version string) (res string, err error) {
	v, err := parse(version)
	if err != nil {
		return
	}

	bumped, err := v.SetPrerelease("")
	res = bumped.String()
	return
}
