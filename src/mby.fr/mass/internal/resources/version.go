package resources

import(
	"fmt"
	"mby.fr/mass/version"
)

var AlreadyBumped = fmt.Errorf("Resource already bumped")
var AlreadyPromoted = fmt.Errorf("Resource already promoted")
var AlreadyReleased = fmt.Errorf("Resource already released")

type Versioner interface {
	Resource
	Version() string
	//FullName() string
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

func writeVersionable(v *Versionable) (err error) {
	var i interface{} = v
	res, ok := i.(Resource)
	if ok {
		err = Write(res)
	}
	return
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
	err = writeVersionable(v)
	if err != nil {
		return
	}

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
	if isDev {
		toVer, err = version.NextRc(fromVer)
		if err != nil {
			return "", err
		}

		v.Ver = toVer
		err = writeVersionable(v)
		if err != nil {
			return
		}
	} else {
		return "", AlreadyPromoted
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
		err = writeVersionable(v)
		if err != nil {
			return
		}
	} else {
		return "", AlreadyReleased
	}

	msg = fmt.Sprintf("%s => %s", fromVer, toVer)
	return
}
