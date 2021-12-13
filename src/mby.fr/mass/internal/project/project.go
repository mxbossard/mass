package project

import (
	//"os"
	//"io/fs"
	"path/filepath"

	"mby.fr/mass/internal/settings"
	"mby.fr/mass/internal/resources"
	//"mby.fr/utils/file"
)

//const DefaultSourceDir = "src"
//const DefaultTestDir = "test"
//const DefaultVersionFile = "version.txt"
//const DefaultInitialVersion = "0.0.1"
//const DefaultBuildFile = "Dockerfile"

var forbiddenNames = []string{resources.DefaultSourceDir, resources.DefaultTestDir, "envs"}

//type Image struct {
//	Name string
//	Dir string
//	SourceDir string
//	TestDir string
//	Buildfile string
//	Version string
//}
//
//type Project struct {
//	Name string
//	Dir string
//	TestDir string
//	Version string
//	Images []Image
//}

func InitProject(name string) (projectPath string, err error) {
	settingsService, err := settings.GetSettingsService()
	if err != nil {
		return
	}

	projectPath = filepath.Join(settingsService.ProjectsDir(), name)

	_, err = resources.Init(projectPath, resources.ProjectKind)
	return
}

//func InitProject0(name string) (projectPath string, err error) {
//	settingsService, err := settings.GetSettingsService()
//	if err != nil {
//		return
//	}
//
//	// Create project dir
//	projectPath, err = file.CreateSubDirectory(settingsService.ProjectsDir(), name)
//	if err != nil {
//		return
//	}
//
//	err = resources.Init(projectPath, resources.ProjectKind)
//	if err != nil {
//		return
//	}
//
//	// Create test dir
//	testDir := testDirpath(projectPath)
//	err = file.CreateDirectory(testDir)
//	if err != nil {
//		return
//	}
//
//	// Init version file
//	versionFile := versionFilepath(projectPath)
//	_, err = file.SoftInitFile(versionFile, resources.DefaultInitialVersion)
//
//	return
//}

func InitImage(projectDir, name string) (imagePath string, err error) {
	imagePath = filepath.Join(projectDir, name)
	_, err = resources.Init(imagePath, resources.ProjectKind)
	return
}

//func InitImage0(projectDir, name string) (imagePath string, err error) {
//	// Create image dir
//	imagePath, err = file.CreateSubDirectory(projectDir, name)
//	if err != nil {
//		return
//	}
//
//	err = resources.Init(imagePath, resources.ProjectKind)
//	if err != nil {
//		return
//	}
//
//	// Create image source dir
//	sourceDir := sourceDirpath(imagePath)
//	err = file.CreateDirectory(sourceDir)
//	if err != nil {
//		return
//	}
//
//	// Create image test dir
//	testDir := testDirpath(imagePath)
//	err = file.CreateDirectory(testDir)
//	if err != nil {
//		return
//	}
//
//	// Init version file
//	versionFile := versionFilepath(imagePath)
//	_, err = file.SoftInitFile(versionFile, resources.DefaultInitialVersion)
//
//	// Init Build file
//	buildfile := buildfileFilepath(imagePath)
//	_, err = file.SoftInitFile(buildfile, "")
//
//	return
//}

//func ListProjects() (projects []Project, err error) {
//	// List directories with a version file to build project list
//        settingsService, err := settings.GetSettingsService()
//        if err != nil {
//                return
//        }
//
//	projectsDir := settingsService.ProjectsDir()
//
//	// Look for a version file then stop walking the branch
//	projectCollector := func(path string, d fs.DirEntry, err error) error {
//		if err != nil {
//			return err
//		}
//		if d.Name() == DefaultVersionFile {
//			// Found a version file
//			parentDir := filepath.Dir(path)
//
//			// => Found a project
//			p, err := buildProject(parentDir)
//			if err != nil {
//				return err
//			}
//			projects = append(projects, p)
//
//			// Stop walking the branch
//			return fs.SkipDir
//		}
//		return nil
//	       }
//
//	filepath.WalkDir(projectsDir, projectCollector)
//
//	return
//}

//func listImages(p Project) (images []Image, err error) {
//        // Look for a Build file then stop walking the branch
//        imageCollector := func(path string, d fs.DirEntry, err error) error {
//                if err != nil {
//                        return err
//                }
//                if d.Name() == DefaultBuildFile {
//                        // Found a version file
//                        parentDir := filepath.Dir(path)
//
//                        // => Found an image
//                        i, err := buildImage(parentDir)
//                        if err != nil {
//                                return err
//                        }
//                        images = append(images, i)
//
//                        // Stop walking the branch
//                        return fs.SkipDir
//                }
//                return nil
//               }
//
//        filepath.WalkDir(p.Dir, imageCollector)
//
//        return
//}
//
//func sourceDirpath(parentPath string) string {
//	return filepath.Join(parentPath, DefaultSourceDir)
//}
//
//func testDirpath(parentPath string) string {
//	return filepath.Join(parentPath, DefaultTestDir)
//}
//
//func versionFilepath(parentPath string) string {
//	return filepath.Join(parentPath, DefaultVersionFile)
//}
//
//func buildfileFilepath(parentPath string) string {
//	return filepath.Join(parentPath, DefaultBuildFile)
//}
//
//func buildProject(path string) (p Project, err error) {
//	path, err = filepath.Abs(path)
//	if err != nil {
//		return
//	}
//	name := filepath.Base(path)
//
//	testDir := testDirpath(path)
//	_, err = os.Stat(testDir); 
//	if os.IsNotExist(err) {
//		testDir = ""
//	} else if err != nil {
//		return
//	}
//
//	versionFile := versionFilepath(path)
//	content, err := os.ReadFile(versionFile); 
//	version := string(content)
//	if os.IsNotExist(err) {
//		version = ""
//	} else if err != nil {
//		return
//	}
//
//	p = Project{Name: name, Dir: path, TestDir: testDir, Version: version}
//	images, err := listImages(p)
//	p.Images = images
//	return
//}
//
//func buildImage(path string) (i Image, err error) {
//	path, err = filepath.Abs(path)
//	if err != nil {
//		return
//	}
//	name := filepath.Base(path)
//
//	sourceDir := sourceDirpath(path)
//	_, err = os.Stat(sourceDir); 
//	if os.IsNotExist(err) {
//		sourceDir = ""
//	} else if err != nil {
//		return
//	}
//
//	testDir := testDirpath(path)
//	_, err = os.Stat(testDir); 
//	if os.IsNotExist(err) {
//		testDir = ""
//	} else if err != nil {
//		return
//	}
//
//	buildfile := buildfileFilepath(path)
//	_, err = os.Stat(buildfile); 
//	if os.IsNotExist(err) {
//		buildfile = ""
//	} else if err != nil {
//		return
//	}
//
//	versionFile := versionFilepath(path)
//	content, err := os.ReadFile(versionFile); 
//	version := string(content)
//	if os.IsNotExist(err) {
//		version = ""
//	} else if err != nil {
//		return
//	}
//
//	i = Image{Name: name, Dir: path, SourceDir: sourceDir, TestDir: testDir, Buildfile: buildfile, Version: version}
//	return
//}

