package resources

import(
	"io/fs"
	"path/filepath"
)

func ScanProjects(path string) (projects []Project, err error) {
	scanner := func(path string, d fs.DirEntry, err error) error {
		if d.Name() == defaultResourceFile {
			parentDir := filepath.Dir(path)
			res, err := Load(parentDir)
			//p, err := buildProject(parentDir)
			//if err != nil {
			//	return err
			//}
			if res.Kind() == ProjectKind {
				projects = append(projects, res)
			}

			return fs.SkipDir
		}
		return nil
	}
	err = filepath.WalkDir(path, scanner)
	if err != nil {
		return
	}
	return
}

func ScanImages(path string) (envs []Env, err error) {
	return
}

func ScanEnvs(path string) (envs []Env, err error) {
	return
}
