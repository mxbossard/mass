package change

import (
	"path/filepath"

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

func calcImageSignature(res resources.Image) (signature string, err error) {
	filesToSign := []string{res.BuildFile, res.SourceDir()}
	signature, err = trust.SignFsContents(filesToSign...)

	// TODO add build config in signature

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

func DoesImageChanged(res resources.Image) (test bool, err error) {
	// Return true if found image changed
	previousSignature, e1 := loadImageSignature(res)
	if e1 != nil {
		return false, e1
	}

	actualSignature, e2 := calcImageSignature(res)
	if e2 != nil {
		return false, e2
	}
	test = previousSignature != actualSignature
	return
}

func calcDeploySignature(res resources.Image) (signature string, err error) {
	// TODO add run config in signature
	// TODO add volumes in signature
	filesToSign := []string{}
	signature, err = trust.SignFsContents(filesToSign...)
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