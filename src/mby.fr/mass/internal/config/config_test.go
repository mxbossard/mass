package config

import (
	"testing"
	"os"

        "github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/require"

	"mby.fr/utils/test"
	"mby.fr/mass/internal/settings"
)

func TestInit(t *testing.T) {
	tempDir, err := test.MkRandTempDir()
	require.NoError(t, err, "should not error")
	require.NoFileExists(t, tempDir, "should not exists")

	// Init Settings for templates to work
	err = settings.Init(tempDir)
	require.NoError(t, err, "should not error")
	os.Chdir(tempDir)

	err = Init(tempDir, nil)
	require.NoError(t, err, "should not error")
	assert.DirExists(t, tempDir, "should exists")
	tempFile := tempDir + "/config.yaml"
	assert.FileExists(t, tempFile, "should exists")
}


func TestRead(t *testing.T) {
	tempDir, err := test.MkRandTempDir()
	require.NoError(t, err, "should not error")
	require.NoFileExists(t, tempDir, "should not exists")
	err = Init(tempDir, nil)
	require.NoError(t, err, "should not error")

	tempFile := tempDir + "/config.yaml"
	c, err := Read(tempFile)
	require.NoError(t, err, "should not error")
	require.NotNil(t, c, "shoult not be nil")
	assert.Len(t, c.Environment, 0, "should be empty")
}

func TestMergeNilStringMaps(t *testing.T) {
	mapA := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
	}
	mergedMap := mergeStringMaps(nil, mapA)
	require.NotNil(t, mergedMap, "should not be nil")
	assert.Equal(t, mapA, mergedMap, "bad merged map")

	mergedMap = mergeStringMaps(mapA, nil)
	require.NotNil(t, mergedMap, "should not be nil")
	assert.Equal(t, mapA, mergedMap, "bad merged map")
}

func TestMergeStringMaps(t *testing.T) {
	mapA := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
	}
	mapB := map[string]string{
		"key1": "valA",
		"key2": "valB",
		"key4": "val4",
	}
	mergedMap := mergeStringMaps(mapA, mapB)
	require.NotNil(t, mergedMap, "should not be nil")
	assert.Len(t, mergedMap, 4, "bad merged map size")
	assert.Equal(t, "valA", mergedMap["key1"], "key not merged")
	assert.Equal(t, "valB", mergedMap["key2"], "key not merged")
	assert.Equal(t, "val3", mergedMap["key3"], "key modified")
	assert.Equal(t, "val4", mergedMap["key4"], "key modified")
}

func TestMergeConfig(t *testing.T) {
	c1 := Config{
		Environment: EnvConfig{
			"key1": "val1",
			"key2": "val2",
		},
	}
	c2 := Config{
		Environment: EnvConfig{
			"key1": "valA",
			"key2": "valB",
			"key3": "val3",
		},
	}
	c3 := Config{
		Environment: EnvConfig{
			"key4": "val4",
			"key1": "valC",
			"key5": "val5",
		},
	}

	mergedConfig := Merge(c1, c2, c3)
	require.NotNil(t, mergedConfig, "should not be nil")
	assert.Len(t, mergedConfig.Environment, 5, "bad merged map size")
	assert.Equal(t, "valC", mergedConfig.Environment["key1"], "key not merged")
	assert.Equal(t, "valB", mergedConfig.Environment["key2"], "key not merged")
	assert.Equal(t, "val3", mergedConfig.Environment["key3"], "key modified")
	assert.Equal(t, "val4", mergedConfig.Environment["key4"], "key modified")
	assert.Equal(t, "val5", mergedConfig.Environment["key5"], "key modified")
}

