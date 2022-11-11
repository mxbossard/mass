package resources

import (
	"fmt"

	"mby.fr/mass/version"
)

var AlreadyBumped = fmt.Errorf("resource already bumped")
var AlreadyPromoted = fmt.Errorf("resource already promoted")
var NotPromoted = fmt.Errorf("resource not promoted yet")
var NotPromotable = fmt.Errorf("resource not promotable")
var AlreadyReleased = fmt.Errorf("resource already released")

type Versioner interface {
	//Resource
	//GetVersionable() *Versionable
	//SetVersionable(Versionable)
	Version() string
	//FullName() string
}

type VersionBumper interface {
	Bump(bool, bool) (string, string, error)
	Promote() (string, string, error)
	Release() (string, string, error)
}

func writeVersionable(v versionable) (err error) {
	err = Write(v.resource)
	return
}

type versionable struct {
	resource Resourcer
	ver string `yaml:"version"`
}

func buildVersionable(res Resourcer, version string) (v versionable) {
	v = versionable{resource: res, ver: version}
	return
}

func (v versionable) init() error {
	return nil
}

func (v versionable) Version() string {
	return v.ver
}

// Bump res version always set qualifier to dev except if qualifier is rc
// Version lifecycle :
// - 1.0.0 -> 1.0.1-dev
// - 1.0.0-rc1 -> 1.0.0-rc2
// - 1.0.3-dev -> 1.0.3-dev
func (v *versionable) Bump(bumpMinor, bumpMajor bool) (toVer, fromVer string, err error) {
	fromVer = v.ver
	var isDev, isRc bool
	if bumpMajor {
		toVer, err = version.NextMajor(fromVer)
		if err != nil {
			return
		}
		toVer, err = version.Dev(toVer)
	} else if bumpMinor {
		toVer, err = version.NextMinor(fromVer)
		if err != nil {
			return
		}
		toVer, err = version.Dev(toVer)
	} else {
		isDev, err = version.IsDev(fromVer)
		if err != nil {
			return
		}
		if isDev {
			err = AlreadyBumped
			return
		}
		isRc, err = version.IsRc(fromVer)
		if err != nil {
			return
		}
		if isRc {
			toVer, err = version.NextRc(fromVer)
		} else {
			toVer, err = version.NextDev(fromVer)
		}
	}
	if err != nil {
		return
	}

	v.ver = toVer
	return
}

// Promote res version from dev to rc.
func (v *versionable) Promote() (toVer, fromVer string, err error) {
	fromVer = v.ver
	var isDev, isRc bool
	isDev, err = version.IsDev(fromVer)
	if err != nil {
		return
	}
	if isDev {
		toVer, err = version.NextRc(fromVer)
		if err != nil {
			return
		}

		v.ver = toVer
		err = writeVersionable(*v)
		return
	}

	isRc, err = version.IsRc(fromVer)
	if err != nil {
		return
	}
	if isRc {
		err = AlreadyPromoted
		return
	}

	err = NotPromotable
	return
}

// Release res version from rc to release
func (v *versionable) Release() (toVer, fromVer string, err error) {
	var isDev, isRc bool
	fromVer = v.ver
	isRc, err = version.IsRc(fromVer)
	if err != nil {
		return
	}
	if isRc {
		toVer, err = version.Release(fromVer)
		if err != nil {
			return
		}

		v.ver = toVer
		err = writeVersionable(*v)
		if err != nil {
			return
		}
	} else {
		isDev, err = version.IsDev(fromVer)
		if err != nil {
			return
		}
		if isDev {
			err = NotPromoted
			return
		}
		err = AlreadyReleased
		return
	}
	return
}
