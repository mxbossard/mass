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

const envDir = "env/"

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
	os.MkdirAll(filepath.Join(path, envDir), 0755)
	for _, e := range envs {
		r, err := BuildEnv("env/" + e)
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
		// --- Resolving relative resource name
		// Project
		{"/", project1, "", project1, ""}, // case 0
		{"/", project1, ProjectKind, project1, ""},
		{project1, project1, ProjectKind, project1, ""},
		{project2, project1, ProjectKind, project1, ""},
		{envDir + env1, project1, ProjectKind, project1, ""},
		{"..", project1, ProjectKind, "", "Unable to found settings path"},
		// Env
		{"/", env1, "", env1, ""},
		{"/", env1, EnvKind, env1, ""},
		{envDir + env1, env1, EnvKind, env1, ""},
		{envDir + env2, env1, EnvKind, env1, ""},
		{project1, env1, EnvKind, env1, ""}, // case 10
		{"..", env1, EnvKind, "", "Unable to found settings path"},
		// Image
		{"/", image11, "", image11, ""},
		{"/", image11, ImageKind, image11, ""},
		{project1, image11, ImageKind, image11, ""},
		{project2, image11, ImageKind, "", ResourceNotFound.Error()},
		{envDir + env1, image11, ImageKind, "", ResourceNotFound.Error()},
		{project1, image11, EnvKind, "", ResourceNotFound.Error()},
		{"..", image11, ImageKind, "", "Unable to found settings path"},

		// Resolving dot resource
		{"/", "", "", "", ResourceNotFound.Error()},
		// Project
		{"/", "", ProjectKind, "", ResourceNotFound.Error()}, // case 20
		{"/", ".", ProjectKind, "", ResourceNotFound.Error()},
		{project1, "", "", project1, ""},
		{project1, ".", "", project1, ""},
		{project1, "", ProjectKind, project1, ""},
		{project1, ".", ProjectKind, project1, ""},
		{envDir + env1, ".", ProjectKind, "", ResourceNotFound.Error()},
		{"..", "", ProjectKind, "", "Unable to found settings path"},
		{"..", ".", ProjectKind, "", "Unable to found settings path"},
		// Env
		{"/", "", EnvKind, "", ResourceNotFound.Error()},
		{"/", ".", EnvKind, "", ResourceNotFound.Error()}, // case 30
		{envDir + env1, "", "", env1, ""},
		{envDir + env1, ".", "", env1, ""},
		{envDir + env1, "", EnvKind, env1, ""},
		{envDir + env1, ".", EnvKind, env1, ""},
		{project1, ".", EnvKind, "", ResourceNotFound.Error()},
		{"..", "", EnvKind, "", "Unable to found settings path"},
		{"..", ".", EnvKind, "", "Unable to found settings path"},
		// Image
		{"/", "", "", "", ResourceNotFound.Error()},
		{"/", ".", ImageKind, "", ResourceNotFound.Error()},
		{project1, "", ImageKind, "", ResourceNotFound.Error()}, // case 40
		{project1, ".", ImageKind, "", ResourceNotFound.Error()},
		{project1 + "/" + image11, "", ImageKind, image11, ""},
		{project1 + "/" + image11, ".", ImageKind, image11, ""},
		{project1 + "/" + image11, "", EnvKind, "", ResourceNotFound.Error()},
		{project1 + "/" + image11, ".", EnvKind, "", ResourceNotFound.Error()},
		{"..", "", ImageKind, "", "Unable to found settings path"},
		{"..", ".", ImageKind, "", "Unable to found settings path"},

		// Resolving absolute resource
		// Project
		{"/", "/" + project1, "", project1, ""},
		{"/", "/" + project1, ProjectKind, project1, ""},
		{project1, "/" + project1, ProjectKind, project1, ""}, // case 50
		{project2, "/" + project1, ProjectKind, project1, ""},
		{envDir + env1, "/" + project1, ProjectKind, project1, ""},
		{"..", "/" + project1, ProjectKind, "", "Unable to found settings path"},
		// Env 
		{"/", "/" + env1, "", env1, ""},
		{"/", "/" + env1, EnvKind, env1, ""},
		{project1, "/" + env1, EnvKind, env1, ""},
		{project2, "/" + env1, EnvKind, env1, ""},
		{envDir + env1, "/" + env1, EnvKind, env1, ""},
		{"..", "/" + env1, EnvKind, "", "Unable to found settings path"},
		// Image 
		{"/", "/" + image11, "", "", ResourceNotFound.Error()}, // case 60
		{"/", "/" + project1 + "/" + image11, "", image11, ""},
		{"/", "/" + image11, ImageKind, "", InconsistentExpression.Error()},
		{"/", "/" + project1 + "/" + image11, ImageKind, image11, ""},
		{project1, "/" + image11, ImageKind, "", InconsistentExpression.Error()},
		{project1, "/" + project1 + "/" + image11, ImageKind, image11, ""},
		{project2, "/" + image11, ImageKind, "", InconsistentExpression.Error()},
		{project2, "/" + project1 + "/" + image11, ImageKind, image11, ""},
		{envDir + env1, "/" + image11, ImageKind, "", InconsistentExpression.Error()},
		{envDir + env1, "/" + project1 + "/" + image11, ImageKind, image11, ""},
		{"..", "/" + image11, ImageKind, "", "Unable to found settings path"}, // case 70

		// Resolving absolute project
		{"/", "project/" + project1, "", project1, ""},
		{"/", "project/" + project1, ProjectKind, project1, ""},
		{"/", "project/" + project1, EnvKind, "", InconsistentExpressionPrefix.Error()},
		{project1, "project/" + project1, ProjectKind, project1, ""},
		{project2, "project/" + project1, ProjectKind, project1, ""},
		{envDir + env1, "project/" + project1, ProjectKind, project1, ""},
		{project2 + "/" + image21, "project/" + project1, ProjectKind, project1, ""},
		{"..", "project/" + project1, ProjectKind, "", "Unable to found settings path"},

		// Resolving absolute env 
		{"/", "env/" + env1, "", env1, ""},
		{"/", "env/" + env1, EnvKind, env1, ""}, // case 80
		{"/", "env/" + env1, ProjectKind, "", InconsistentExpressionPrefix.Error()},
		{project1, "env/" + env1, EnvKind, env1, ""},
		{project2, "env/" + env1, EnvKind, env1, ""},
		{envDir + env1, "env/" + env2, EnvKind, env2, ""},
		{project2 + "/" + image21, "env/" + env1, EnvKind, env1, ""},
		{"..", "env/" + env1, EnvKind, "", "Unable to found settings path"},

		// Resolving absolute image
		{"/", "image/" + project1 + "/" + image11, "", image11, ""},
		{project1, "image/" + project1 + "/" + image11, ImageKind, image11, ""},
		{project1, "image/" + project1 + "/" + image11, ProjectKind, "", InconsistentExpressionPrefix.Error()},
		{project1, "image/" + image11, ImageKind, "", InconsistentExpression.Error()}, // case 90
		{project2, "image/" + project1 + "/" + image11, ImageKind, image11, ""},
		{envDir + env1, "image/" + project1 + "/" + image11, ImageKind, image11, ""},
		{project2 + "/" + image21, "image/" + project1 + "/" + image11, ImageKind, image11, ""},
		{"..", "image/" + project1 + "/" + image11, ImageKind, "", "Unable to found settings path"},

		// Resolving not existing resources
		{"/", "notExisting", ProjectKind, "", ResourceNotFound.Error()},
		{"/", "/notExisting", ProjectKind, "", ResourceNotFound.Error()},
		{"/", "project/notExisting", ProjectKind, "", ResourceNotFound.Error()},
		{"/", "env/notExisting", EnvKind, "", ResourceNotFound.Error()},
		{"/", "image/notExisting/notExisting", ImageKind, "", ResourceNotFound.Error()},
		{project1, "notExisting", ProjectKind, "", ResourceNotFound.Error()},
		{project1, "/notExisting", ProjectKind, "", ResourceNotFound.Error()},
		{project1, "project/notExisting", ProjectKind, "", ResourceNotFound.Error()},
		{project1, "env/notExisting", EnvKind, "", ResourceNotFound.Error()},
		{project1, "image/notExisting/notExisting", ImageKind, "", ResourceNotFound.Error()},
		{"..", "notExisting", ProjectKind, "", "Unable to found settings path"},
		{"..", "/notExisting", ProjectKind, "", "Unable to found settings path"},
		{"..", "project/notExisting", ProjectKind, "", "Unable to found settings path"},

		// Resolving bad kind resources
		{"/", project1, EnvKind, "", ResourceNotFound.Error()},
		{"/", "/" + project1, EnvKind, "", ResourceNotFound.Error()},
		{"/", "project/" + project1, EnvKind, "", InconsistentExpressionPrefix.Error()},
		{project1, project1, EnvKind, "", ResourceNotFound.Error()},
		{project1, "/" + project1, EnvKind, "", ResourceNotFound.Error()},
		{project1, "project/" + project1, EnvKind, "", InconsistentExpressionPrefix.Error()},
		{"/", "env/" + project1, ProjectKind, "", InconsistentExpressionPrefix.Error()},
		{"/", "env/" + project1, EnvKind, "", ResourceNotFound.Error()},
		{"/", "env/" + project1, ImageKind, "", InconsistentExpressionPrefix.Error()},

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
