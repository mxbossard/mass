package resources

import(
	"testing"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/require"

	"mby.fr/utils/test"
	"mby.fr/mass/internal/settings"
)

const env1 = "env1"
const env2 = "env2"
const env3 = "env3"

const project1 = "p1"
const project2 = "p2"
const project3 = "p3"

const image11 = "i11"
const image12 = "i12"
const image13 = "i13"

const image21 = "i21"
const image22 = "i22"
const image23 = "i23"

const image31 = "i31"
const image32 = "i32"
const image33 = "i33"

var (
	envs = []string{env1, env2, env3}
	projects = []string{project1, project2, project3}
	images = map[string][]string{
		project1: []string{image11, image12, image13},
		project2: []string{image21, image22, image33},
		project3: []string{image31, image32, image33},
	}
)
func initWorkspace(t *testing.T) (path string) {
	// Build fake workspace with resources tree
	path, err := test.MkRandTempDir()
	require.NoError(t, err, "should not error")

	// Init Settings for templates to work
        err = settings.Init(path)
        require.NoError(t, err, "should not error")
        os.Chdir(path)

	// Init envs
	for _, e := range envs {
		r, err := BuildEnv(e)
		require.NoError(t, err, "should not error")
		err = r.Init()
		require.NoError(t, err, "should not error")
	}

	// Init projects
	for _, p := range projects {
		r, err := BuildProject(p)
		require.NoError(t, err, "should not error")
		err = r.Init()
		require.NoError(t, err, "should not error")

		// Init project images
		for _, i := range images[p] {
			r, err := BuildImage(p + "/" + i)
			require.NoError(t, err, "should not error")
			err = r.Init()
			require.NoError(t, err, "should not error")
		}
	}

	return
}

func TestResolveResource(t *testing.T) {
	fakeWorkspacePath := initWorkspace(t)

	cases := []struct {
		fromPath, exprIn, kindIn, resNameWanted, errWanted string
	} {
		{".", project1, ProjectKind, project1, ""},
		{project1, project1, ProjectKind, project1, ""},
		{project2, project1, ProjectKind, project1, ""},
		{project1, project1, EnvKind, "", "Resource not found"},
		{"..", project1, ProjectKind, "", "Unable to found settings path"},

		{".", "/" + project1, ProjectKind, project1, ""},
		{project1, "/" + project1, ProjectKind, project1, ""},
		{project2, "/" + project1, ProjectKind, project1, ""},
		{project1, "/" + project1, EnvKind, "", "Resource not found"},
		{"..", "/" + project1, ProjectKind, "", "Unable to found settings path"},

		{".", "project/" + project1, ProjectKind, project1, ""},
		{project1, "project/" + project1, ProjectKind, project1, ""},
		{project2, "project/" + project1, ProjectKind, project1, ""},
		{project2, "project/" + project1, EnvKind, "", "Resource not found"},
		{"..", "project/" + project1, ProjectKind, "", "Unable to found settings path"},

	}

	for i, c := range cases {
		path := filepath.Join(fakeWorkspacePath, c.fromPath)
		err := os.Chdir(path)
		require.NoError(t, err, "should not error for chdir on case %d", i)
		r, err := resolveResource(c.exprIn, c.kindIn)
		if c.errWanted == "" {
			require.NoError(t, err, "should not error on case %d", i)
		} else {
			assert.EqualError(t, err, c.errWanted, "bad error for case %d", i)
		}

		if c.resNameWanted == "" {
			assert.Nil(t, r, "should not found a resource for case %d", i)
		} else {
			require.NotNil(t, r, "should found a resource for case %d", i)
			assert.Equal(t, c.resNameWanted, r.Name(), "Bad resource name for case %d", i)
		}
	}
}
