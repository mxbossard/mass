package workspace

import (
	"fmt"
	//"strings"
	"path/filepath"

	//"mby.fr/utils/filez"
	"mby.fr/mass/internal/resources"
	"mby.fr/mass/internal/settings"
)

type ResourceDontExists struct {
	kind resources.Kind
	name string
}

func (e ResourceDontExists) Error() string {
	return fmt.Sprintf("%s: %s don't exists", e.kind, e.name)
}

// Assert resource exists by name. Use complete name for images.
func AssertResourceExists(kind resources.Kind, name string) (err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return
	}

	switch kind {
	case resources.EnvKind:
		parentDir := ss.EnvsDir()
		resourceDir := filepath.Join(parentDir, name)
		err = resources.AssertResourceDir(kind, resourceDir)
	case resources.ProjectKind:
		parentDir := ss.ProjectsDir()
		resourceDir := filepath.Join(parentDir, name)
		err = resources.AssertResourceDir(kind, resourceDir)
	case resources.ImageKind:
		parentDir := ss.ProjectsDir()
		resourceDir := filepath.Join(parentDir, name)
		err = resources.AssertResourceDir(kind, resourceDir)
	}

	if err != nil {
		return ResourceDontExists{kind, name}
	}
	return nil
}
