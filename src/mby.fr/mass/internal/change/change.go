package change

import (
	"io/fs"
	"path/filepath"
	"strings"

	"mby.fr/mass/internal/resources"
	"mby.fr/mass/internal/settings"
	"mby.fr/utils/cache"
	"mby.fr/utils/trust"
)

const defaultImageCacheDir = "imageSignatures"
const defaultDeployCacheDir = "deploySignatures"

var imageCacheDir cache.Cache
var deployCacheDir cache.Cache

func Init() (err error) {
	if imageCacheDir != nil && deployCacheDir != nil {
		return
	}

	// Initializes caches
	ss, e := settings.GetSettingsService()
	if err != nil {
		return e
	}
	imageSignaturesCacheDir := filepath.Join(ss.CacheDir(), defaultImageCacheDir)
	deploySignaturesCacheDir := filepath.Join(ss.CacheDir(), defaultDeployCacheDir)

	imageCacheDir, err = cache.NewPersistentCache(imageSignaturesCacheDir)
	if err != nil {
		return
	}
	deployCacheDir, err = cache.NewPersistentCache(deploySignaturesCacheDir)
	if err != nil {
		return
	}

	return
}

func fileTree(rootPath string) (tree string, err error) {
	files := []string{}
	fn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		entry := ""
		if d.IsDir() {
			entry += "d"
		} else {
			entry += "-"
		}
		entry += path
		files = append(files, entry)
		return nil
	}
	err = filepath.WalkDir(rootPath, fn)
	if err != nil {
		return "", err
	}
	tree = strings.Join(files, ";")
	return
}

func calcImageSignature(res resources.Image) (signature string, err error) {
	filesToSign := []string{res.BuildFile, res.AbsSourceDir()}
	filesSignature, err := trust.SignFsContents(filesToSign...)
	if err != nil {
		return "", err
	}

	configs, err := resources.MergedConfig(res)
	if err != nil {
		return "", err
	}
	fileTree, err := fileTree(res.AbsSourceDir())
	if err != nil {
		return "", err
	}

	signature, err = trust.SignObjects(configs.BuildArgs, filesSignature, fileTree)

	return
}

func imageCacheKey(res resources.Image) (signature string) {
	return res.FullName()
}

func loadImageSignature(res resources.Image) (signature string, err error) {
	key := imageCacheKey(res)
	value, ok, e := imageCacheDir.LoadString(key)
	if e != nil {
		return signature, e
	}
	if ok {
		signature = value
	}
	return
}

func StoreImageSignature(res resources.Image) (err error) {
	signature, e := calcImageSignature(res)
	if e != nil {
		return e
	}
	key := imageCacheKey(res)
	err = imageCacheDir.StoreString(key, signature)
	return
}

func DoesImageChanged(res resources.Image) (changed bool, sign string, err error) {
	// Return true if found image changed
	previousSignature, e1 := loadImageSignature(res)
	if e1 != nil {
		return false, "", e1
	}

	actualSignature, e2 := calcImageSignature(res)
	if e2 != nil {
		return false, "", e2
	}
	changed = previousSignature != actualSignature
	sign = actualSignature
	return
}

func calcDeploySignature(res resources.Image) (signature string, err error) {
	// TODO add volumes in signature
	filesToSign := []string{}
	filesSignature, err := trust.SignFsContents(filesToSign...)

	configs, err := resources.MergedConfig(res)
	if err != nil {
		return "", err
	}
	signature, err = trust.SignObjects(filesSignature, configs.Environment, configs.RunArgs)

	return
}

func deployCacheKey(res resources.Image) (signature string) {
	return res.FullName()
}

func loadDeploySignature(res resources.Image) (signature string, err error) {
	key := deployCacheKey(res)
	value, ok, e := deployCacheDir.LoadString(key)
	if e != nil {
		return signature, e
	}
	if ok {
		signature = value
	}
	return
}

func StoreDeploySignature(res resources.Image) (err error) {
	signature, e := calcDeploySignature(res)
	if e != nil {
		return e
	}
	key := deployCacheKey(res)
	err = deployCacheDir.StoreString(key, signature)
	return
}

func DoesDeployChanged(res resources.Image) (test bool, err error) {
	// Return true if found deploy changed
	previousSignature, e1 := loadDeploySignature(res)
	if e1 != nil {
		return false, e1
	}

	actualSignature, e2 := calcDeploySignature(res)
	if e2 != nil {
		return false, e2
	}
	test = previousSignature != actualSignature
	return
}
