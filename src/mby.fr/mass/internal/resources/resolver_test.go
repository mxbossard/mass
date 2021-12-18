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
		// Resolving relative resource name
		{"/", project1, "", project1, ""}, // case 0
		{"/", project1, ProjectKind, project1, ""}, 
		{project1, project1, ProjectKind, project1, ""},
		{project2, project1, ProjectKind, project1, ""},
		{"..", project1, ProjectKind, "", "Unable to found settings path"},
		{"/", env1, "", env1, ""}, // Case 5
		{"/", env1, EnvKind, env1, ""},
		{env1, env1, EnvKind, env1, ""},
		{env2, env1, EnvKind, env1, ""},
		{"..", env1, EnvKind, "", "Unable to found settings path"},

		// Resolving dot resource
		{"/", "", ProjectKind, "", ResourceNotFound.Error()}, // Case 5
		{"/", ".", ProjectKind, "", ResourceNotFound.Error()},
		{project1, "", ProjectKind, project1, ""},
		{project1, ".", ProjectKind, project1, ""},
		{"..", "", ProjectKind, "", "Unable to found settings path"},
		{"..", ".", ProjectKind, "", "Unable to found settings path"}, // Case 10

		// Resolving absolute resource
		{"/", "/" + project1, ProjectKind, project1, ""},
		{project1, "/" + project1, ProjectKind, project1, ""},
		{project2, "/" + project1, ProjectKind, project1, ""},
		{project2, "/" + project1, ProjectKind, project1, ""}, // duplicate
		{"..", "/" + project1, ProjectKind, "", "Unable to found settings path"}, // Case 15

		// Resolving absolute project
		{"/", "project/" + project1, ProjectKind, project1, ""},
		{project1, "project/" + project1, ProjectKind, project1, ""},
		{project2, "project/" + project1, ProjectKind, project1, ""},
		{project2, "project/" + project1, ProjectKind, project1, ""}, // duplicate
		{"..", "project/" + project1, ProjectKind, "", "Unable to found settings path"}, // Case 20

		// Resolving not existing resources
		{"/", "notExisting", ProjectKind, "", ResourceNotFound.Error()},
		{"/", "/notExisting", ProjectKind, "", ResourceNotFound.Error()},
		{"/", "project/notExisting", ProjectKind, "", ResourceNotFound.Error()},
		{project1, "notExisting", ProjectKind, "", ResourceNotFound.Error()},
		{project1, "/notExisting", ProjectKind, "", ResourceNotFound.Error()}, // case 25
		{project1, "project/notExisting", ProjectKind, "", ResourceNotFound.Error()},
		{"..", "notExisting", ProjectKind, "", "Unable to found settings path"},
		{"..", "/notExisting", ProjectKind, "", "Unable to found settings path"},
		{"..", "project/notExisting", ProjectKind, "", "Unable to found settings path"},

		// Resolving bad kind resources
		{"/", project1, EnvKind, "", ResourceNotFound.Error()}, // case 30
		{"/", "/" + project1, EnvKind, "", ResourceNotFound.Error()},
		{"/", "project/" + project1, EnvKind, "", InconsistentResourceKind.Error()},
		{project1, project1, EnvKind, "", ResourceNotFound.Error()},
		{project1, "/" + project1, EnvKind, "", ResourceNotFound.Error()},
		{project1, "project/" + project1, EnvKind, "", InconsistentResourceKind.Error()},
		{"/", "env/" + project1, ProjectKind, "", InconsistentResourceKind.Error()},
		{"/", "env/" + project1, EnvKind, "", ResourceNotFound.Error()},
		{"/", "env/" + project1, ImageKind, "", InconsistentResourceKind.Error()},

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
