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

func TestRead(t *testing.T) {
	expectedRendering := "foo: {{ .Bar }}."
	data, err := read(testTemplate)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedRendering, data, "bad reading")

	// re-read
	data, err = read(testTemplate)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedRendering, data, "bad reading")
}

func TestReadNewline(t *testing.T) {
	expectedRendering := "foo \nbar\nbaz\n"
	data, err := read(testNewlineTemplate)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedRendering, data, "bad reading")
}

func TestRender(t *testing.T) {
	builder := strings.Builder{}
	err := Render(testTemplate, &builder, nil)
	require.Error(t, err, "should error")

	builder.Reset()
	barValue := "baz"
	expectedRendering := "foo: " + barValue + "."
	data := struct{ Bar string } { Bar: barValue }
	err = Render(testTemplate, &builder, data)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedRendering, builder.String(), "bad rendering")

}

func TestRenderToFile(t *testing.T) {
	tempFile, err := test.BuildRandTempPath()
	require.NoFileExists(t, tempFile, "should not exists")

	err = RenderToFile(testTemplate, tempFile, nil)
	require.Error(t, err, "should error")
	require.NoFileExists(t, tempFile, "should not exists")

	barValue := "baz"
	expectedRendering := "foo: " + barValue + "."
	data := struct{ Bar string } { Bar: barValue }
	err = RenderToFile(testTemplate, tempFile, data)
	require.NoError(t, err, "should not error")
	require.FileExists(t, tempFile, "should exists")
	content, err := os.ReadFile(tempFile)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedRendering, string(content), "bad rendering")
}

