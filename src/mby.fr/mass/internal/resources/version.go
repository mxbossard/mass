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
	Version() string
	//FullName() string
}

type VersionBumper interface {
	Bump(bool, bool) (string, error)
	Promote() (string, error)
	Release() (string, error)
}

type Versionable struct {
	Ver string `yaml:"version"`
}

func (v Versionable) Version() string {
	return v.Ver
}

// Bump res version always set qualifier to dev except if qualifier is rc
// Version lifecycle :
// - 1.0.0 -> 1.0.1-dev
// - 1.0.0-rc1 -> 1.0.0-rc2
// - 1.0.3-dev -> 1.0.3-dev
func (v *Versionable) Bump(bumpMinor, bumpMajor bool) (msg string, err error) {
	var fromVer, toVer string
	fromVer = v.Version()
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
		isDev, err := version.IsDev(fromVer)
		if err != nil {
			return "", err
		}
		if isDev {
			return "", AlreadyBumped
		}
		isRc, err := version.IsRc(fromVer)
		if err != nil {
			return "", err
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

	v.Ver = toVer

	msg = fmt.Sprintf("%s => %s", fromVer, toVer)
	return
}

// Promote res version from dev to rc.
func (v *Versionable) Promote() (msg string, err error) {
	var fromVer, toVer string
	fromVer = v.Version()
	isDev, err := version.IsDev(fromVer)
	if err != nil {
		return "", err
	}
	isRc, err := version.IsRc(fromVer)
	if err != nil {
		return "", err
	}
	if isDev {
		toVer, err = version.NextRc(fromVer)
		if err != nil {
			return "", err
		}

		v.Ver = toVer
	} else if isRc {
		return "", AlreadyPromoted
	} else {
		return "", NotPromotable
	}

	msg = fmt.Sprintf("%s => %s", fromVer, toVer)
	return
}

// Release res version from rc to release
func (v *Versionable) Release() (msg string, err error) {
	var fromVer, toVer string
	fromVer = v.Version()
	isRc, err := version.IsRc(fromVer)
	if err != nil {
		return "", err
	}
	if isRc {
		toVer, err = version.Release(fromVer)
		if err != nil {
			return "", err
		}

		v.Ver = toVer
	} else {
		isDev, err := version.IsDev(fromVer)
		if err != nil {
			return "", err
		}
		if isDev {
			return "", NotPromoted
		}
		return "", AlreadyReleased
	}

	msg = fmt.Sprintf("%s => %s", fromVer, toVer)
	return
}

func Bump(r Resource, minor, major bool) (msg string, err error) {
	switch v := r.(type) {
	case VersionBumper:
		msg, err = v.Bump(minor, major)
		if err != nil {
			return
		}
		err = Write(r)
	default:
		err = fmt.Errorf("not able to bump resource of type: %T", r)
	}
	return
}

func Promote(r Resource) (msg string, err error) {
	switch v := r.(type) {
	case VersionBumper:
		msg, err = v.Promote()
		if err != nil {
			return
		}
		err = Write(r)
	default:
		err = fmt.Errorf("not able to promote resource of type: %T", r)
	}
	return
}

func Release(r Resource) (msg string, err error) {
	switch v := r.(type) {
	case VersionBumper:
		msg, err = v.Release()
		if err != nil {
			return
		}
		err = Write(r)
	default:
		err = fmt.Errorf("not able to release resource of type: %T", r)
	}
	return
}
