package resources

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypeFromKind(t *testing.T) {
	typ := TypeFromKind(EnvKind)
	assert.Equal(t, reflect.TypeOf((*Env)(nil)).Elem(), typ, "bad type for env kind")
}

func TestKindFromType(t *testing.T) {
	typ := TypeFromKind(EnvKind)
	kind := KindFromType(typ)
	assert.Equal(t, EnvKind, kind, "bad kind for Env type")
}

func TestKindFromResource(t *testing.T) {
	var env Env
	kind := KindFromResource(env)
	assert.Equal(t, EnvKind, kind, "bad kind for Env resource")
}

func TestEmptyKindSet(t *testing.T) {
	s := *NewKindSet()
	//assert.Len(t, s, 0, "should contains no kind")
	assert.Len(t, s, 1, "should contains AllKind")
	assert.Contains(t, s, AllKind, "should contains AllKind")
}

func TestKindSet(t *testing.T) {
	s := *NewKindSet(AllKind)
	assert.Len(t, s, 1, "should contains AllKind")
	assert.Contains(t, s, AllKind, "should contains AllKind")

	s = *NewKindSet(EnvKind)
	assert.Len(t, s, 1, "should contains EnvKind")
	assert.Contains(t, s, EnvKind, "should contains EnvKind")

	s = *NewKindSet(EnvKind, EnvKind)
	assert.Len(t, s, 1, "should contains EnvKind")
	assert.Contains(t, s, EnvKind, "should contains EnvKind")

	s = *NewKindSet(EnvKind, AllKind)
	assert.Len(t, s, 1, "should contains only AllKind")
	assert.Contains(t, s, AllKind, "should contains AllKind")

	s = *NewKindSet(ProjectKind, EnvKind)
	assert.Len(t, s, 2, "should contains EnvKind")
	assert.Contains(t, s, EnvKind, "should contains EnvKind")
	assert.Contains(t, s, ProjectKind, "should contains ProjectKind")
}

func TestKindSetString(t *testing.T) {
	s := *NewKindSet(AllKind)
	assert.Equal(t, "all", s.String())

	s = *NewKindSet(EnvKind)
	assert.Equal(t, "env", s.String())

	s = *NewKindSet(EnvKind, ProjectKind)
	assert.Equal(t, "env,project", s.String())

	s = *NewKindSet(EnvKind, ProjectKind, ImageKind)
	assert.Equal(t, "env,image,project", s.String())
}

func TestKindFromAlias(t *testing.T) {
	k, ok := KindFromAlias("")
	assert.False(t, ok)

	k, ok = KindFromAlias("foo")
	assert.False(t, ok)

	k, ok = KindFromAlias("all")
	assert.True(t, ok)
	assert.Equal(t, AllKind, k)

	k, ok = KindFromAlias("envs")
	assert.True(t, ok)
	assert.Equal(t, EnvKind, k)
}

func TestKindSetContains(t *testing.T) {
	s := *NewKindSet(EnvKind, ProjectKind)
	test := s.Contains(AllKind)
	assert.False(t, test)

	test = s.Contains(ImageKind)
	assert.False(t, test)

	test = s.Contains(ProjectKind)
	assert.True(t, test)

	test = s.Contains(ImageKind)
	assert.False(t, test)

	s = *NewKindSet(AllKind)
	test = s.Contains(AllKind)
	assert.True(t, test)

	test = s.Contains(ImageKind)
	assert.True(t, test)

	test = s.Contains(ProjectKind)
	assert.True(t, test)

	test = s.Contains(ImageKind)
	assert.True(t, test)
}

func TestMarshalYAML(t *testing.T) {
	yaml, err := EnvKind.MarshalYAML()
	require.NoError(t, err, "should not error")
	assert.Equal(t, EnvKind.String(), yaml)
}

func TestUnmarshalYAML(t *testing.T) {
	var k Kind
	unmarshal := func(i interface{}) error {
		s := EnvKind.String()
		i = &s
		return nil
	}
	err := k.UnmarshalYAML(unmarshal)
	require.NoError(t, err, "should not error")
	assert.Equal(t, EnvKind, k, "bad kind")
}
