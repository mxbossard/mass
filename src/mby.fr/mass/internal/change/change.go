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

func loadImageSignature(res resources.Resource) (err error) {
	// Use utils/cache to load signature
	return
}

func StoreImageSignature(res resources.Resource) (err error) {
	// Use utils/cache to store signature
	return
}

func DoesImageChanged(res resources.Resource) (res bool, err error) {
	// Return true if found image changed

	// Need to compare versus a store if the image change since last build.
	return
}

func StoreDeploySignature(res resources.Resource) (err error) {
	return
}

func DoesDeployChanged(res resources.Resource) (res bool, err error) {
	return
}
