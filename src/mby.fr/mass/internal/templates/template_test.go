package templates

import (
	"testing"
	"strings"

        "github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/require"
)

const testTemplate = "test.txt"

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

func TestRender(t *testing.T) {
	builder := strings.Builder{}
	//expectedRendering := "foo: <no value>."
	err := Render(testTemplate, &builder, nil)
	require.Error(t, err, "should error")
	//assert.Equal(t, expectedRendering, builder.String(), "bad rendering")

	builder.Reset()
	barValue := "baz"
	expectedRendering := "foo: " + barValue + "."
	data := struct{ Bar string } { Bar: barValue }
	err = Render(testTemplate, &builder, data)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedRendering, builder.String(), "bad rendering")

}
