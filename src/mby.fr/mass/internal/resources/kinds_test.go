package resources

import(
	"testing"

	"github.com/stretchr/testify/assert"
        //"github.com/stretchr/testify/require"
)


func TestEmptyKinds(t *testing.T) {
	s := Kinds()
	assert.Len(t, s, 0, "should contains no kind")
	//assert.Len(t, s, 1, "should contains AllKind")
	//assert.Equal(t, AllKind, s[0], "should contains AllKind")
}

func TestKinds(t *testing.T) {
	s := Kinds(AllKind)
	assert.Len(t, s, 1, "should contains AllKind")
	assert.Contains(t, s, AllKind, "should contains AllKind")

	s = Kinds(EnvKind)
	assert.Len(t, s, 1, "should contains EnvKind")
	assert.Contains(t, s, EnvKind, "should contains EnvKind")

	s = Kinds(EnvKind, EnvKind)
	assert.Len(t, s, 1, "should contains EnvKind")
	assert.Contains(t, s, EnvKind, "should contains EnvKind")

	s = Kinds(EnvKind, AllKind)
	assert.Len(t, s, 1, "should contains only AllKind")
	assert.Contains(t, s, AllKind, "should contains AllKind")

	s = Kinds(ProjectKind, EnvKind)
	assert.Len(t, s, 2, "should contains EnvKind")
	assert.Contains(t, s, EnvKind, "should contains EnvKind")
	assert.Contains(t, s, ProjectKind, "should contains ProjectKind")

}

