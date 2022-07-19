package change

import (
	"mby.fr/mass/internal/settings"
	"mby.fr/utils/cache"
)

const defaultImageCacheDir = "imageSignatures"
const defaultDeployCacheDir = "deploySignatures"

var imageCacheDir cache.Cache
var deployCacheDir cache.Cache

func init() (err error) {
	// Initializes caches
	var ss, e := settings.GetSettingsService()
	if err != nil {
		return e
	}
	var imageSignaturesCacheDir := filepath.Join(ss.CacheDir(), defaultImageCacheDir)
	var deploySignaturesCacheDir := filepath.Join(ss.CacheDir(), defaultDeployCacheDir)

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

func calcImageSignature(res resources.Resource) (signature string, err error) {
	// Use utils/trust to calc signature
	return
}

func imageCacheKey(res resources.Resource) (signature string) {
	return
}

func loadImageSignature(res resources.Resource) ((signature string, err error) {
	var key := imageCacheKey(res)
	value, ok, e := imageCacheDir.LoadString(key)
	if e != nil {
		return _, e
	}
	return
}

func StoreImageSignature(res resources.Resource) (err error) {
	var signature, e := calcImageSignature(res)
	if e != nil {
		return e
	}
	var key := imageCacheKey(res)
	err = imageCacheDir.StoreString(key, signature)
	return
}

func DoesImageChanged(res resources.Resource) (res bool, err error) {
	// Return true if found image changed

	previousSignature, e1 := loadImageSignature(res)
	if e1 != nil {
		return e1
	}

	actualSignature, e2 := calcImageSignature(res)
	if e2 != nil {
		return e2
	}

	return previousSignature != actualSignature
}

func StoreDeploySignature(res resources.Resource) (err error) {
	return
}

func DoesDeployChanged(res resources.Resource) (res bool, err error) {
	return
}
