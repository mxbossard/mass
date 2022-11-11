package resources

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

const PathSeparator string = string(filepath.Separator)

var MaxDepthError error = fmt.Errorf("max depth reached")

func pathDepth(path string) (depth int) {
	cleanedPath := filepath.Clean(path)
	return strings.Count(cleanedPath, PathSeparator)
}

func buildScanner0(resKind Kind, c chan<- interface{}) fs.WalkDirFunc {
	scanner := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		//fmt.Println("scanning", path)
		if d.Name() == DefaultResourceFile {
			parentDir := filepath.Dir(path)
			res, err := ReadResourcer(parentDir)
			if err != nil {
				return err
			}
			if res.Kind() == resKind {
				c <- res
				return fs.SkipDir
			}
		}
		return nil
	}
	return scanner
}

func buildScanner(rootPath string, resKind Kind, maxDepth int, c chan<- interface{}) fs.WalkDirFunc {
	if maxDepth >= 0 {
		rootPathDepth := pathDepth(rootPath)
		maxDepth += rootPathDepth
	}
	scanner := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if maxDepth >= 0 && pathDepth(path) > maxDepth+1 {
			//fmt.Printf("Reached max depth: %d with path: %s\n", maxDepth, path)
			return fs.SkipDir
		}

		//fmt.Println("scanning", path)
		if d.Name() == DefaultResourceFile {
			parentDir := filepath.Dir(path)
			res, err := ReadResourcer(parentDir)
			if err != nil {
				return err
			}
			if res.Kind() == resKind || resKind == AllKind {
				c <- res
				return fs.SkipDir
			}
		}
		return nil
	}
	return scanner
}

func buildScanner2[T Resourcer](rootPath string, maxDepth int, c chan<- T) fs.WalkDirFunc {
	if maxDepth >= 0 {
		rootPathDepth := pathDepth(rootPath)
		maxDepth += rootPathDepth
	}
	scanner := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if maxDepth >= 0 && pathDepth(path) > maxDepth+1 {
			//fmt.Printf("Reached max depth: %d with path: %s\n", maxDepth, path)
			return fs.SkipDir
		}

		//fmt.Println("scanning", path)
		if d.Name() == DefaultResourceFile {
			parentDir := filepath.Dir(path)
			res, err := Read[T](parentDir)
			if IsBadResourceType(err) {
				// pass we are scanning
			} else if err != nil {
				return err
			} else {
				c <- res
				return fs.SkipDir
			}
		}
		return nil
	}
	return scanner
}

func ScanMaxDepth[T Resourcer](path string, maxDepth int) (resources []T, err error) {
	c := make(chan T)
	finished := make(chan bool)
	go func() {
		// consume not buffered channel in a goroutine to avoid to be stuck
		for r := range c {
			//fmt.Println("Consuming project")
			resources = append(resources, r)
		}
		finished <- true
		close(finished)
	}()

	scanner := buildScanner2[T](path, maxDepth, c)
	err = filepath.WalkDir(path, scanner)
	close(c)
	if errors.Is(err, fs.ErrNotExist) {
		// Swallow error if path don't exists
		err = nil
		return
	} else if err != nil {
		return
	}
	// BLock until array finished
	<-finished
	return
}

func Scan[T Resourcer](path string) (projects []T, err error) {
	return ScanMaxDepth[T](path, -1)
}

func scanResourcesFrom(fromDir string, resourceKind Kind) (resources []Resourcer, err error) {
	c := make(chan interface{})
	finished := make(chan bool)
	go func() {
		// consume not buffered channel in a goroutine to avoid to be stuck
		for r := range c {
			//fmt.Println("Consuming project")
			resources = append(resources, r.(Resourcer))
		}
		finished <- true
		close(finished)
	}()

	scanner := buildScanner(fromDir, resourceKind, 1, c)
	err = filepath.WalkDir(fromDir, scanner)
	close(c)
	if errors.Is(err, fs.ErrNotExist) {
		// Swallow error if path don't exists
		err = nil
		return
	} else if err != nil {
		return
	}
	// BLock until array finished
	<-finished
	return
}

/*
func ScanProjectsMaxDepth(path string, maxDepth int) (projects []Project, err error) {
	//c := make(chan interface{})
	c := make(chan Project)
	finished := make(chan bool)
	go func() {
		// consume not buffered channel in a goroutine to avoid to be stuck
		for r := range c {
			//fmt.Println("Consuming project")
			projects = append(projects, r)
		}
		finished <- true
		close(finished)
	}()

	//scanner := buildScanner(path, ProjectKind, maxDepth, c)
	scanner := buildScanner2[Project](path, maxDepth, c)
	err = filepath.WalkDir(path, scanner)
	close(c)
	if errors.Is(err, fs.ErrNotExist) {
		// Swallow error if path don't exists
		err = nil
		return
	} else if err != nil {
		return
	}
	// BLock until array finished
	<-finished
	return
}

func ScanProjects(path string) (projects []Project, err error) {
	return ScanProjectsMaxDepth(path, -1)
}

func ScanImagesMaxDepth(path string, maxDepth int) (images []*Image, err error) {
	c := make(chan interface{})
	finished := make(chan bool)
	go func() {
		// consume not buffered channel in a goroutine to avoid to be stuck
		for r := range c {
			//fmt.Println("Consuming image")
			images = append(images, r.(*Image))
		}
		finished <- true
		close(finished)
	}()

	scanner := buildScanner(path, ImageKind, maxDepth, c)
	err = filepath.WalkDir(path, scanner)
	close(c)
	if errors.Is(err, fs.ErrNotExist) {
		// Swallow error if path don't exists
		err = nil
		return
	} else if err != nil {
		return
	}
	// BLock until array finished
	<-finished
	return
}

func ScanImages(path string) (images []*Image, err error) {
	return ScanImagesMaxDepth(path, -1)
}

func ScanEnvsMaxDepth(path string, maxDepth int) (envs []*Env, err error) {
	c := make(chan interface{})
	finished := make(chan bool)
	go func() {
		// consume not buffered channel in a goroutine to avoid to be stuck
		for r := range c {
			//fmt.Println("Consuming env")
			envs = append(envs, r.(*Env))
		}
		finished <- true
		close(finished)
	}()

	scanner := buildScanner(path, EnvKind, maxDepth, c)
	err = filepath.WalkDir(path, scanner)
	close(c)
	if errors.Is(err, fs.ErrNotExist) || errors.Is(err, MaxDepthError) {
		// Swallow error if path don't exists
		err = nil
		return
	} else if err != nil {
		return
	}
	// BLock until array finished
	<-finished
	return
}

func ScanEnvs(path string) (envs []*Env, err error) {
	return ScanEnvsMaxDepth(path, -1)
}

func scanResourcesFrom(fromDir string, resourceKind Kind) (res []Resourcer, err error) {
	var envs []*Env
	var projects []Project
	var images []*Image
	switch resourceKind {
	case AllKind:
		envs, err = ScanEnvsMaxDepth(fromDir, 1)
		if err != nil && !IsResourceNotFound(err) {
			return
		}
		projects, err = ScanProjectsMaxDepth(fromDir, 1)
		if err != nil && !IsResourceNotFound(err) {
			return
		}
		images, err = ScanImagesMaxDepth(fromDir, 1)
		if err != nil && !IsResourceNotFound(err) {
			return
		}

	case EnvKind:
		envs, err = ScanEnvsMaxDepth(fromDir, -1)
		if err != nil {
			return
		}
	case ProjectKind:
		projects, err = ScanProjectsMaxDepth(fromDir, -1)
		if err != nil {
			return
		}
	case ImageKind:
		images, err = ScanImagesMaxDepth(fromDir, -1)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("Not supported kind")
		return
	}

	for _, r := range envs {
		res = append(res, r)
	}
	for _, r := range projects {
		res = append(res, r)
	}
	for _, r := range images {
		res = append(res, r)
	}

	if IsResourceNotFound(err) && len(res) > 0 {
		// Swallow ResourceNotFound error if found something
		err = nil
	}

	return
}
*/