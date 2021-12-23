package resources

import(
	"fmt"
	"strings"
	//"path/filepath"

	//"mby.fr/utils/file"
)

type BadResourceName struct {
	kind Kind
	name string
}
func (e BadResourceName) Error() string {
	return fmt.Sprintf("bad %s name: %s", e.kind, e.name)
}

type NotResourceDir struct {
	kind Kind
	dir string
}
func (e NotResourceDir) Error() string {
	return fmt.Sprintf("directory: %s is not a %s", e.dir, e.kind)
}

func AssertResourceName(kind Kind, name string) error {
	switch kind {
	case EnvKind, ProjectKind:
		if strings.Contains(name, "/") {
			return BadResourceName{kind, name}
		}
	case ImageKind:
		if strings.Count(name, "/") > 1 {
			return BadResourceName{kind, name}
		}
	}
	return nil
}

func AssertResourceDir(kind Kind, dir string) error {
	r, err := Read(dir)
	if err != nil {
		return err
	}
	if r.Kind() != kind {
		return NotResourceDir{kind, dir}
	}
	return nil
}

