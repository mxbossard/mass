package resources

import(
	//"fmt"
	"io/fs"
	"path/filepath"
)

func buildScanner(resKind Kind, c chan<- interface{}) (fs.WalkDirFunc) {
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

func ScanProjects(path string) (projects []Project, err error) {
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
	scanner := buildScanner(ProjectKind, c)
	err = filepath.WalkDir(path, scanner)
	close(c)
	if err != nil {
		return
	}
	// BLock until array finished
	<- finished
	return
}

func ScanImages(path string) (images []Image, err error) {
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
	scanner := buildScanner(ImageKind, c)
	err = filepath.WalkDir(path, scanner)
	close(c)
	if err != nil {
		return
	}
	// BLock until array finished
	<- finished
	return
}

func ScanEnvs(path string) (envs []Env, err error) {
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
	scanner := buildScanner(EnvKind, c)
	err = filepath.WalkDir(path, scanner)
	close(c)
	if err != nil {
		return
	}
	// BLock until array finished
	<- finished
	return
}
