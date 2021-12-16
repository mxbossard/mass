package templates

import (
	"testing"
	"strings"
	"os"

        "github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/require"

	"mby.fr/utils/test"
)

const testTemplate = "test.txt"
const testNewlineTemplate = "testNewline.txt"

func TestInit(t *testing.T) {
	tempDir, err := test.MkRandTempDir()
	require.NoError(t, err, "should not error")
	err = Init(tempDir)
	require.NoError(t, err, "should not error")
	require.DirExists(t, tempDir + "/src", "should be copied")
	require.FileExists(t, tempDir + "/src/" + testTemplate, "should be copied")
	require.FileExists(t, tempDir + "/src/" + testNewlineTemplate, "should be copied")
}

func TestInitWithNotExistingDir(t *testing.T) {
	tempFile, err := test.BuildRandTempPath()
	require.NoError(t, err, "should not error")
	err = Init(tempFile)
	require.Error(t, err, "should error")
}

func TestReadFromEmbeded(t *testing.T) {
	templatesDir := ""
	err := Init(templatesDir)
	require.NoError(t, err, "should not error")
	
	r := New(templatesDir)
	expectedRendering := "foo: {{ .Bar }}.\n"
	data, err := r.read(testTemplate)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedRendering, data, "bad reading")

	// re-read
	data, err = r.read(testTemplate)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedRendering, data, "bad reading")
}

func TestReadFromDir(t *testing.T) {
	tempDir, err := test.MkRandTempDir()
	require.NoError(t, err, "should not error")
	err = Init(tempDir)
	require.NoError(t, err, "should not error")

	r := New(tempDir)
	expectedRendering := "foo: {{ .Bar }}.\n"
	data, err := r.read(testTemplate)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedRendering, data, "bad reading")

	// re-read
	data, err = r.read(testTemplate)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedRendering, data, "bad reading")
}

func TestReadNewlineFromEmbeded(t *testing.T) {
	r := New("")

	expectedRendering := "foo \nbar\nbaz\n"
	data, err := r.read(testNewlineTemplate)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedRendering, data, "bad reading")
}

func TestReadNewlineFromDir(t *testing.T) {
	tempDir, err := test.MkRandTempDir()
	require.NoError(t, err, "should not error")
	err = Init(tempDir)
	require.NoError(t, err, "should not error")

	r := New(tempDir)

	expectedRendering := "foo \nbar\nbaz\n"
	data, err := r.read(testNewlineTemplate)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedRendering, data, "bad reading")
}

func TestRender(t *testing.T) {
	r := New("")

	builder := strings.Builder{}
	err := r.Render(testTemplate, &builder, nil)
	require.Error(t, err, "should error")

	builder.Reset()
	barValue := "baz"
	expectedRendering := "foo: " + barValue + ".\n"
	data := struct{ Bar string } { Bar: barValue }
	err = r.Render(testTemplate, &builder, data)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedRendering, builder.String(), "bad rendering")

}

func TestRenderToFile(t *testing.T) {
	r := New("")

	tempFile, err := test.BuildRandTempPath()
	require.NoFileExists(t, tempFile, "should not exists")

	err = r.RenderToFile(testTemplate, tempFile, nil)
	require.Error(t, err, "should error")
	require.NoFileExists(t, tempFile, "should not exists")

	barValue := "baz"
	expectedRendering := "foo: " + barValue + ".\n"
	data := struct{ Bar string } { Bar: barValue }
	err = r.RenderToFile(testTemplate, tempFile, data)
	require.NoError(t, err, "should not error")
	require.FileExists(t, tempFile, "should exists")
	content, err := os.ReadFile(tempFile)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedRendering, string(content), "bad rendering")
}

