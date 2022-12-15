package resources

import (
	"os"
	"path/filepath"
	"strings"

	"mby.fr/mass/internal/settings"
	"mby.fr/utils/filez"
)

type Image struct {
	directoryBase  	`yaml:"base,inline"`
	configurableDir `yaml:"-"` // Ignore this field for yaml marshalling
	testable    	`yaml:"testable,inline"`
	versionable 	`yaml:"versionable,inline"`

	SourceDirectory string  `yaml:"sourceDirectory"`
	BuildFile       string  `yaml:"buildFile"`
	Project         Project `yaml:"-"` // Ignore this field for yaml marshalling
}

func (i Image) init() (err error) {
	err = i.directoryBase.init()
	if err != nil {
		return
	}
	err = i.testable.init()
	if err != nil {
		return
	}
	err = i.versionable.init()
	if err != nil {
		return
	}

	// Init source directory
	err = os.MkdirAll(i.AbsSourceDir(), 0755)

	// Init Build file
	buildfileContent := ""
	//buildfileContent := "FROM alpine\n"
	_, err = filez.SoftInitFile(i.BuildFile, buildfileContent)

	return
}

func (i Image) AbsSourceDir() string {
	return absResourcePath(i.Dir(), i.SourceDirectory)
}

func (i Image) AbsBuildFile() string {
	return absResourcePath(i.Dir(), i.BuildFile)
}

func (i Image) FullName() string {
	return i.Project.FullName() + "/" + i.ImageName()
}

func (i Image) ImageName() string {
	return i.directoryBase.name
}

func (i Image) FullImageName() string {
	if i.Version() != "" {
		return strings.ToLower(i.FullName()) + ":" + i.Version()
	} else {
		return strings.ToLower(i.FullName()) + ":latest"
	}
}

func (i Image) AbsoluteName() (name string, err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return "", err
	}
	name = ss.Settings().Name + "-" + i.FullName()
	return
}

func (i Image) Match(name string, k Kind) bool {
	return i.directoryBase.Match(name, k) || name == i.ImageName() && KindsMatch(k, i.Kind())
}

func buildImage(projectDir, imageName string) (r Image, err error) {
	imageDir := filepath.Join(projectDir, imageName)
	version := DefaultInitialVersion
	buildfile := DefaultBuildFile
	sourceDir := DefaultSourceDir

	project, err := Read[Project](projectDir) //buildProject(projectPath)
	if err != nil {
		return
	}

	b, err := buildDirectoryBase(ImageKind, imageDir)
	if err != nil {
		return
	}

	r = Image{
		directoryBase:   b,
		BuildFile:       buildfile,
		SourceDirectory: sourceDir,
		Project:         project,
	}
	r.configurableDir = buildConfigurableDir(b)

	t, err := buildTestable(r)
	if err != nil {
		return
	}
	r.testable = t

	versionable := buildVersionable(r, version)
	r.versionable = versionable

	return
}