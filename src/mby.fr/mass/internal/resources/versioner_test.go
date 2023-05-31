package resources

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"mby.fr/mass/internal/commontest"
)

func TestType(t *testing.T) {
	wksPath := commontest.InitTempWorkspace(t)
	defer os.RemoveAll(wksPath)
	p, err := Build[Project]("p1", wksPath)
	require.NoError(t, err, "should not error")
	image, err := Build[Image]("i1", p)
	require.NoError(t, err, "should not error")

	assert.Implements(t, (*VersionBumper)(nil), &image, "image pointer should implements Versioner")
}

func TestBump(t *testing.T) {
	wksPath := commontest.InitTempWorkspace(t)
	defer os.RemoveAll(wksPath)
	p, err := Build[Project]("p1", wksPath)
	require.NoError(t, err, "should not error")
	image, err := Build[Image]("i1", p)
	require.NoError(t, err, "should not error")

	assert.Equal(t, "0.1.0-dev", image.Version(), "Bad initial version")

	toVer, fromVer, err := image.Bump(false, false)
	require.Error(t, err, "Bump must return an error")
	assert.Equal(t, AlreadyBumped, err, "Bad bump error")
	assert.Equal(t, "", toVer, "Bad bumped message")
	assert.Equal(t, image.Version(), fromVer, "Bad bumped message")

	image.versionable.ver = "2.0.1-rc3"
	toVer, fromVer, err = image.Bump(false, false)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "2.0.1-rc4", image.Version(), "Bad bumped version")
	assert.Equal(t, "2.0.1-rc4", toVer, "Bad bumped message")
	assert.Equal(t, "2.0.1-rc3", fromVer, "Bad bumped message")

	image.versionable.ver = "3.0.3"
	toVer, fromVer, err = image.Bump(false, false)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "3.0.4-dev", image.Version(), "Bad bumped version")
	assert.Equal(t, "3.0.4-dev", toVer, "Bad bumped message")
	assert.Equal(t, "3.0.3", fromVer, "Bad bumped message")

	// Bump next major
	image.versionable.ver = "4.0.1-rc3"
	toVer, fromVer, err = image.Bump(false, true)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "5.0.0-dev", image.Version(), "Bad bumped version")
	assert.Equal(t, "5.0.0-dev", toVer, "Bad bumped message")
	assert.Equal(t, "4.0.1-rc3", fromVer, "Bad bumped message")

	image.versionable.ver = "5.0.3"
	toVer, fromVer, err = image.Bump(false, true)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "6.0.0-dev", image.Version(), "Bad bumped version")
	assert.Equal(t, "6.0.0-dev", toVer, "Bad bumped message")
	assert.Equal(t, "5.0.3", fromVer, "Bad bumped message")

	image.versionable.ver = "4.0.1-rc3"
	toVer, fromVer, err = image.Bump(false, true)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "5.0.0-dev", image.Version(), "Bad bumped version")
	assert.Equal(t, "5.0.0-dev", toVer, "Bad bumped message")
	assert.Equal(t, "4.0.1-rc3", fromVer, "Bad bumped message")

	image.versionable.ver = "5.0.3"
	toVer, fromVer, err = image.Bump(false, true)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "6.0.0-dev", image.Version(), "Bad bumped version")
	assert.Equal(t, "6.0.0-dev", toVer, "Bad bumped message")
	assert.Equal(t, "5.0.3", fromVer, "Bad bumped message")

	// Bump next minor
	image.versionable.ver = "6.0.1-rc3"
	toVer, fromVer, err = image.Bump(false, true)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "7.0.0-dev", image.Version(), "Bad bumped version")
	assert.Equal(t, "7.0.0-dev", toVer, "Bad bumped message")
	assert.Equal(t, "6.0.1-rc3", fromVer, "Bad bumped message")

	image.versionable.ver = "7.0.3"
	toVer, fromVer, err = image.Bump(false, true)
	require.NoError(t, err, "Bump must not return an error")
	assert.Equal(t, "8.0.0-dev", image.Version(), "Bad bumped version")
	assert.Equal(t, "8.0.0-dev", toVer, "Bad bumped message")
	assert.Equal(t, "7.0.3", fromVer, "Bad bumped message")
}

func TestPromote(t *testing.T) {
	wksPath := commontest.InitTempWorkspace(t)
	defer os.RemoveAll(wksPath)
	p, err := Build[Project]("p1", wksPath)
	require.NoError(t, err, "should not error")
	image, err := Build[Image]("i1", p)
	require.NoError(t, err, "should not error")

	toVer, fromVer, err := image.Promote()
	require.NoError(t, err, "must not return an error")
	assert.Equal(t, "0.1.0-rc1", image.Version(), "Bad promoted version")
	assert.Equal(t, "0.1.0-rc1", toVer, "Bad bumped message")
	assert.Equal(t, "0.1.0-dev", fromVer, "Bad bumped message")

	image.versionable.ver = "2.0.1-rc3"
	toVer, fromVer, err = image.Promote()
	require.Error(t, err, "must return an error")
	assert.Equal(t, AlreadyPromoted, err, "Bad promote error")
	assert.Equal(t, "", toVer, "Bad bumped message")
	assert.Equal(t, image.Version(), fromVer, "Bad bumped message")

	image.versionable.ver = "3.0.3"
	toVer, fromVer, err = image.Promote()
	require.Error(t, err, "must return an error")
	assert.Equal(t, NotPromotable, err, "Bad promote error")
	assert.Equal(t, "", toVer, "Bad bumped message")
	assert.Equal(t, image.Version(), fromVer, "Bad bumped message")
}

func TestRelease(t *testing.T) {
	wksPath := commontest.InitTempWorkspace(t)
	defer os.RemoveAll(wksPath)
	p, err := Build[Project]("p1", wksPath)
	require.NoError(t, err, "should not error")
	image, err := Build[Image]("i1", p)
	require.NoError(t, err, "should not error")

	toVer, fromVer, err := image.Release()
	require.Error(t, err, "must return an error")
	assert.Equal(t, NotPromoted, err, "Bad release error")
	assert.Equal(t, "", toVer, "Bad bumped message")
	assert.Equal(t, image.Version(), fromVer, "Bad bumped message")

	image.versionable.ver = "2.0.1-rc3"
	toVer, fromVer, err = image.Release()
	require.NoError(t, err, "must not return an error")
	assert.Equal(t, "2.0.1", image.Version(), "Bad released version")
	assert.Equal(t, "2.0.1", toVer, "Bad bumped message")
	assert.Equal(t, "2.0.1-rc3", fromVer, "Bad bumped message")

	image.versionable.ver = "3.0.3"
	toVer, fromVer, err = image.Release()
	require.Error(t, err, "must return an error")
	assert.Equal(t, AlreadyReleased, err, "Bad release error")
	assert.Equal(t, "", toVer, "Bad bumped message")
	assert.Equal(t, image.Version(), fromVer, "Bad bumped message")
}
