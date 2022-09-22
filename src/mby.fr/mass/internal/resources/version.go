package resources

import(
	"fmt"
	"mby.fr/mass/version"
)

var AlreadyBumped = fmt.Errorf("resource already bumped")
var AlreadyPromoted = fmt.Errorf("resource already promoted")
var NotPromoted = fmt.Errorf("resource not promoted yet")
var AlreadyReleased = fmt.Errorf("resource already released")

type Versioner interface {
	//Resource
	//GetVersionable() *Versionable
	//SetVersionable(Versionable)
	Version() string
	FullName() string
	Bump(bool, bool) (string, error)
	Promote() (string, error)
	Release() (string, error)
	
}

func writeVersionable(v *Versionable) (err error) {
	var i interface{} = *v
	res, ok := i.(Resource)
	if ok {
		err = Write(res)
	} else {
		err = fmt.Errorf("unable to write Versionable of type %T", v)
	}
	return
}

func writeVersioner(v Versioner) (err error) {
	res, ok := v.(Resource)
	if ok {
		err = Write(res)
	} else {
		err = fmt.Errorf("unable to write Versioner of type %T", v)
	}
	return
}

type Versionable struct {
	ver string `yaml:"version"`
}

func (v Versionable) version() string {
	return v.ver
}

// Bump res version always set qualifier to dev except if qualifier is rc
// Version lifecycle : 
// - 1.0.0 -> 1.0.1-dev
// - 1.0.0-rc1 -> 1.0.0-rc2
// - 1.0.3-dev -> 1.0.3-dev
func (v *Versionable) bump(bumpMinor, bumpMajor bool) (msg string, err error) {
	var fromVer, toVer string
	fromVer = v.ver
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

	v.ver = toVer

	msg = fmt.Sprintf("%s => %s", fromVer, toVer)
	return
}

// Promote res version from dev to rc.
func (v *Versionable) promote() (msg string, err error) {
	var fromVer, toVer string
	fromVer = v.ver
	isDev, err := version.IsDev(fromVer) 
	if err != nil {
		return "", err
	}
	if isDev {
		toVer, err = version.NextRc(fromVer)
		if err != nil {
			return "", err
		}

		v.ver = toVer
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
func (v *Versionable) release() (msg string, err error) {
	var fromVer, toVer string
	fromVer = v.ver
	isRc, err := version.IsRc(fromVer) 
	if err != nil {
		return "", err
	}
	if isRc {
		toVer, err = version.Release(fromVer)
		if err != nil {
			return "", err
		}

		v.ver = toVer
		err = writeVersionable(v)
		if err != nil {
			return
		}
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

/*
func Bump(v Versioner, bumpMinor, bumpMajor bool) (msg string, err error) {
	versionablePtr := v.GetVersionable()
	//versionablePtr.Ver = toVer
	msg, err = versionablePtr.Bump(bumpMinor, bumpMajor)
	if err != nil {
		return
	}

	v.SetVersionable(*versionablePtr)

	err = writeVersioner(v)
	if err != nil {
		return
	}

	return
}
*/