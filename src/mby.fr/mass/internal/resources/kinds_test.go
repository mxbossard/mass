package resources

import(
	"testing"

	"github.com/stretchr/testify/assert"
        //"github.com/stretchr/testify/require"
)


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

	s = *NewKindSet(EnvKind,ProjectKind)
	assert.Equal(t, "env,project", s.String())

	s = *NewKindSet(EnvKind,ProjectKind,ImageKind)
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
	s := *NewKindSet(EnvKind,ProjectKind)
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
