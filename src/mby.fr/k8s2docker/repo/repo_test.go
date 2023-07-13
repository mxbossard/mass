package repo

import (
	_ "io"
	_ "net/http"
	_ "net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testDbPath = "/tmp/mydb"
)

func init() {
}

func TestReadWrite_String(t *testing.T) {
	initDb(testDbPath)
	defer os.RemoveAll(testDbPath)

	key := "key"
	expectedValue := "foo"

	var v string
	var err error

	// Read from not existing collection
	v, err = read[string]("strings", key)
	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)

	// Write
	err = write("strings", key, expectedValue)
	require.NoError(t, err)

	// Read existing
	v, err = read[string]("strings", key)
	require.NoError(t, err)
	assert.Equal(t, expectedValue, v)

	// Read not existing
	v, err = read[string]("strings", "otherKey")
	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestReadWrite_Map(t *testing.T) {
	initDb(testDbPath)
	defer os.RemoveAll(testDbPath)

	key := "key"
	expectedValue := map[string]any{
		"a": "foo",
		"b": "bar",
		"c": []any{
			map[string]any{"k1": "v1", "k2": float64(1)},
			map[string]any{"k1": "v2", "k2": float64(2)},
		},
		"d": map[string]any{"p1": float64(3), "p2": false},
	}

	var v map[string]any
	var err error

	// Read from not existing collection
	v, err = read[map[string]any]("maps", key)
	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)

	// Write
	err = write("maps", key, expectedValue)
	require.NoError(t, err)

	// Read existing
	v, err = read[map[string]any]("maps", key)
	require.NoError(t, err)
	assert.Equal(t, expectedValue, v)

	// Read not existing
	v, err = read[map[string]any]("maps", "otherKey")
	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestReadWrite_Struct(t *testing.T) {
	//TODO
}
