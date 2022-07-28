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
			res, err := Read(parentDir)
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

func buildScanner(resKind Kind, maxDepth int, c chan<- interface{}) fs.WalkDirFunc {
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
			res, err := Read(parentDir)
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

func ScanProjectsMaxDepth(path string, maxDepth int) (projects []Project, err error) {
	c := make(chan interface{})
	finished := make(chan bool)
	go func() {
		// consume not buffered channel in a goroutine to avoid to be stuck
		for r := range c {
			//fmt.Println("Consuming project")
			projects = append(projects, r.(Project))
		}
		finished <- true
		close(finished)
	}()
	scanner := buildScanner(ProjectKind, maxDepth, c)
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

func ScanImagesMaxDepth(path string, maxDepth int) (images []Image, err error) {
	c := make(chan interface{})
	finished := make(chan bool)
	go func() {
		// consume not buffered channel in a goroutine to avoid to be stuck
		for r := range c {
			//fmt.Println("Consuming image")
			images = append(images, r.(Image))
		}
		finished <- true
		close(finished)
	}()
	scanner := buildScanner(ImageKind, maxDepth, c)
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

func ScanImages(path string) (images []Image, err error) {
	return ScanImagesMaxDepth(path, -1)
}

func ScanEnvsMaxDepth(path string, maxDepth int) (envs []Env, err error) {
	c := make(chan interface{})
	finished := make(chan bool)
	go func() {
		// consume not buffered channel in a goroutine to avoid to be stuck
		for r := range c {
			//fmt.Println("Consuming env")
			envs = append(envs, r.(Env))
		}
		finished <- true
		close(finished)
	}()
	scanner := buildScanner(EnvKind, maxDepth, c)
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

func ScanEnvs(path string) (envs []Env, err error) {
	return ScanEnvsMaxDepth(path, -1)
}
