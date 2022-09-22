package resources

import (
	//"fmt"
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/mass/internal/commontest"
)

func TestType(t *testing.T) {
	wksPath := commontest.InitTempWorkspace(t)
	defer os.RemoveAll(wksPath)
	_, imagePath := commontest.InitRandImage(t, wksPath)
	image, err := BuildImage(imagePath)
	require.NoError(t, err, "should not return an error")

	assert.Implements(t, (*Versioner)(nil), &image, "image pointer should implements Versioner")
}

func TestBump(t *testing.T) {
	wksPath := commontest.InitTempWorkspace(t)
	defer os.RemoveAll(wksPath)
	_, imagePath := commontest.InitRandImage(t, wksPath)
	image, err := BuildImage(imagePath)
	require.NoError(t, err, "should not return an error")
	assert.Equal(t, "0.0.1-dev", image.Version(), "Bad initial version")

	msg, err := image.Bump(false, false)
	//msg, err := Bump(&image, false, false)
	require.Error(t, err, "Bump must return an error")
	assert.Equal(t, AlreadyBumped, err, "Bad bump error")
	assert.Equal(t, "", msg, "Bad bumped message")

	image.Versionable.ver = "2.0.1-rc3"
	msg, err = image.Bump(false, false)
	//msg, err = Bump(&image, false, false)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "2.0.1-rc4", image.Version(), "Bad bumped version")
	assert.Equal(t, "2.0.1-rc3 => 2.0.1-rc4", msg, "Bad bumped message")

	image.Versionable.ver = "3.0.3"
	msg, err = image.Bump(false, false)
	//msg, err = Bump(&image, false, false)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "3.0.4-dev", image.Version(), "Bad bumped version")
	assert.Equal(t, "3.0.3 => 3.0.4-dev", msg, "Bad bumped message")

	// Bump next major
	image.Versionable.ver = "4.0.1-rc3"
	msg, err = image.Bump(false, true)
	//msg, err = Bump(&image, false, true)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "5.0.0-dev", image.Version(), "Bad bumped version")
	assert.Equal(t, "4.0.1-rc3 => 5.0.0-dev", msg, "Bad bumped message")

	image.Versionable.ver = "5.0.3"
	msg, err = image.Bump(false, true)
	//msg, err = Bump(&image, false, true)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "6.0.0-dev", image.Version(), "Bad bumped version")
	assert.Equal(t, "5.0.3 => 6.0.0-dev", msg, "Bad bumped message")

	image.Versionable.ver = "4.0.1-rc3"
	msg, err = image.Bump(false, true)
	//msg, err = Bump(&image, false, true)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "5.0.0-dev", image.Version(), "Bad bumped version")
	assert.Equal(t, "4.0.1-rc3 => 5.0.0-dev", msg, "Bad bumped message")

	image.Versionable.ver = "5.0.3"
	msg, err = image.Bump(false, true)
	//msg, err = Bump(&image, false, true)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "6.0.0-dev", image.Version(), "Bad bumped version")
	assert.Equal(t, "5.0.3 => 6.0.0-dev", msg, "Bad bumped message")

	// Bump next minor
	image.Versionable.ver = "6.0.1-rc3"
	msg, err = image.Bump(false, true)
	//msg, err = Bump(&image, true, false)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "6.1.0-dev", image.Version(), "Bad bumped version")
	assert.Equal(t, "6.0.1-rc3 => 6.1.0-dev", msg, "Bad bumped message")

	image.Versionable.ver = "7.0.3"
	msg, err = image.Bump(false, true)
	//msg, err = Bump(&image, true, false)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "7.1.0-dev", image.Version(), "Bad bumped version")
	assert.Equal(t, "7.0.3 => 7.1.0-dev", msg, "Bad bumped message")
}

func TestPromote(t *testing.T) {
	wksPath := commontest.InitTempWorkspace(t)
	defer os.RemoveAll(wksPath)
	_, imagePath := commontest.InitRandImage(t, wksPath)
	image, err := BuildImage(imagePath)
	require.NoError(t, err, "should not return an error")

	msg, err := image.Promote()
	require.NoError(t, err, "must not return an error")
	assert.Equal(t, "0.0.1-rc1", image.Version(), "Bad promoted version")
	assert.Equal(t, "0.0.1-dev => 0.0.1-rc1", msg, "Bad promoted message")

	image.Versionable.ver = "2.0.1-rc3"
	msg, err = image.Promote()
	require.Error(t, err, "must return an error")
	assert.Equal(t, AlreadyPromoted, err, "Bad promote error")
	assert.Equal(t, "", msg, "Bad promote message")

	image.Versionable.ver = "3.0.3"
	msg, err = image.Promote()
	require.Error(t, err, "must return an error")
	assert.Equal(t, AlreadyPromoted, err, "Bad promote error")
	assert.Equal(t, "", msg, "Bad promote message")
}

func TestRelease(t *testing.T) {
	wksPath := commontest.InitTempWorkspace(t)
	defer os.RemoveAll(wksPath)
	_, imagePath := commontest.InitRandImage(t, wksPath)
	image, err := BuildImage(imagePath)
	require.NoError(t, err, "should not return an error")

	msg, err := image.Release()
	require.Error(t, err, "must return an error")
	assert.Equal(t, NotPromoted, err, "Bad release error")
	assert.Equal(t, "", msg, "Bad release message")

	image.Versionable.ver = "2.0.1-rc3"
	msg, err = image.Release()
	require.NoError(t, err, "must not return an error")
	assert.Equal(t, "2.0.1", image.Version(), "Bad released version")
	assert.Equal(t, "2.0.1-rc3 => 2.0.1", msg, "Bad released message")

	image.Versionable.ver = "3.0.3"
	msg, err = image.Release()
	require.Error(t, err, "must return an error")
	assert.Equal(t, AlreadyReleased, err, "Bad release error")
	assert.Equal(t, "", msg, "Bad release message")
}
