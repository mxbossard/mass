package resources

import(
	"testing"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/require"

	"mby.fr/utils/test"
	//"mby.fr/utils/errorz"
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

func TestSplitExpression(t *testing.T) {
	cases := []struct {
                exprIn string
                kindWanted Kind
                nameWanted string
        } {
		{"foo", AllKind, "foo"},
		{"p/bar", ProjectKind, "bar"},
		{"project/bar", ProjectKind, "bar"},
		{"projects/bar", ProjectKind, "bar"},
		{"foo/bar", AllKind, "foo/bar"},
		{"e/bar", EnvKind, "bar"},
		{"env/bar", EnvKind, "bar"},
		{"envs/bar", EnvKind, "bar"},
		{"i/bar", ImageKind, "bar"},
		{"i/bar/baz", ImageKind, "bar/baz"},
		{"image/bar", ImageKind, "bar"},
		{"images/bar", ImageKind, "bar"},
	}
	for i, c := range cases {
		kind, name := splitExpression(c.exprIn)
		assert.Equal(t, c.kindWanted, kind, "bad kind for case %d", i)
		assert.Equal(t, c.nameWanted, name, "bad name for case %d", i)
	}

}

func TestSplitImageName(t *testing.T) {
	cases := []struct {
                nameIn string
                projectWanted string
                imageWanted string
        } {
		{"foo", "", "foo"},
		{"foo/bar", "foo", "bar"},
		{"", "", ""},
		{"foo/bar/baz", "foo", "bar/baz"},
	}
	for i, c := range cases {
		project, image := splitImageName(c.nameIn)
		assert.Equal(t, c.projectWanted, project, "bad project name for case %d", i)
		assert.Equal(t, c.imageWanted, image, "bad image name for case %d", i)
	}

}

func TestResolveResourceFrom(t *testing.T) {
	fakeWorkspacePath := initWorkspace(t)

	cases := []struct {
		fromPath, exprIn string
		kindIn Kind
		resNameWanted string
		errWanted error
	} {
		// --- Resolving relative resource name
		// Project
		{fakeWorkspacePath, project1, "", "", InvalidArgument}, // case 0
		{fakeWorkspacePath, "", ProjectKind, "", InvalidArgument},
		{"", project1, ProjectKind, "", InvalidArgument},
		{fakeWorkspacePath, project1, ProjectKind, project1, nil},
		{fakeWorkspacePath, project1, EnvKind, "", ResourceNotFound},
		{fakeWorkspacePath, project1, ImageKind, "", ResourceNotFound},
		{fakeWorkspacePath + "/" + project1, project1, ProjectKind, project1, nil},
		{fakeWorkspacePath + "/" + project2, project1, ProjectKind, "", ResourceNotFound},
		{fakeWorkspacePath + "/" + envDir + env1, project1, ProjectKind, "", ResourceNotFound},

		// Env
		{fakeWorkspacePath, env1, "", "", InvalidArgument},
		{fakeWorkspacePath, env1, EnvKind, env1, nil}, // case 10
		{fakeWorkspacePath, env1, ProjectKind, "", ResourceNotFound},
		{fakeWorkspacePath, env1, ImageKind, "", ResourceNotFound},
		{fakeWorkspacePath + "/" + envDir + "/" + env1, env1, EnvKind, env1, nil},
		{fakeWorkspacePath + "/" + envDir + "/" + env2, env1, EnvKind, "", ResourceNotFound},
		{fakeWorkspacePath + "/" + project2, env1, EnvKind, "", ResourceNotFound},

		// Image
		{fakeWorkspacePath, image11, "", "", InvalidArgument},
		{fakeWorkspacePath, image11, ImageKind, "", ResourceNotFound},
		{fakeWorkspacePath, project1 + "/" + image11, ImageKind, project1 + "/" + image11, nil},
		{fakeWorkspacePath, project1 + "/" + image11, ProjectKind, "", ResourceNotFound},
		{fakeWorkspacePath, project1 + "/" + image11, EnvKind, "", ResourceNotFound}, // case 20
		{fakeWorkspacePath, project2 + "/" + image11, ImageKind, "", ResourceNotFound},
		{fakeWorkspacePath + "/" + project1, project1 + "/" + image11, ImageKind, project1 + "/" + image11, nil},
		{fakeWorkspacePath + "/" + project1 + "/" + image11, project1 + "/" + image11, ImageKind, project1 + "/" + image11, nil},
		{fakeWorkspacePath + "/" + project2, project1 + "/" + image11, ImageKind, "", ResourceNotFound},
		{fakeWorkspacePath + "/" + project2 + "/" + image21, project1 + "/" + image11, ImageKind, "", ResourceNotFound},

		// Resolving not existing resources
		{fakeWorkspacePath, "notExisting", ProjectKind, "", ResourceNotFound},
		{fakeWorkspacePath, "notExisting", EnvKind, "", ResourceNotFound},
		{fakeWorkspacePath, "notExisting", ImageKind, "", ResourceNotFound}, // case 28
	}

	for i, c := range cases {
		r, err := resolveResourceFrom(c.fromPath, c.exprIn, c.kindIn)
		if c.errWanted == nil {
			assert.NoError(t, err, "should not error on case %d", i)
		} else {
			assert.ErrorIs(t, err, c.errWanted, "bad error for case %d", i)
		}

		if c.resNameWanted == "" {
			assert.Nil(t, r, "should not found a resource for case %d", i)
		} else {
			require.NotNil(t, r, "should found a resource for case %d", i)
			assert.Equal(t, c.resNameWanted, r.Name(), "Bad resource name for case %d", i)
		}
	}
}

func TestResolveContextualResource(t *testing.T) {
	fakeWorkspacePath := initWorkspace(t)

	cases := []struct {
		fromPath, exprIn string
		kindIn Kind
		resNameWanted string
		errWanted error
	} {
		// Project
		{"/", project1, "", "", InvalidArgument}, // case 0
		{"/", "", ProjectKind, "", InvalidArgument},
		{"/", project1, ProjectKind, project1, nil},
		{"/", project1, EnvKind, "", ResourceNotFound},
		{"/", project1, ImageKind, "", ResourceNotFound},
		{"/" + project1, project1, ProjectKind, project1, nil},
		{"/" + project2, project1, ProjectKind, project1, nil},
		{"/" + envDir + env1, project1, ProjectKind, project1, nil},

		// Env
		{"/", env1, "", "", InvalidArgument},
		{"/", "", EnvKind, "", InvalidArgument},
		{"/", env1, EnvKind, env1, nil}, // case 10
		{"/", env1, ProjectKind, "", ResourceNotFound},
		{"/", env1, ImageKind, "", ResourceNotFound},
		{"/" + envDir + "/" + env1, env1, EnvKind, env1, nil},
		{"/" + envDir + "/" + env2, env1, EnvKind, env1, nil},
		{"/" + project2, env1, EnvKind, env1, nil},

		// Image absolute
		{"/", image11, "", "", InvalidArgument},
		{"/", "", ImageKind, "", InvalidArgument},
		{"/", image11, ImageKind, "", ResourceNotFound},
		{"/", project1 + "/" + image11, ImageKind, project1 + "/" + image11, nil},
		{"/", project1 + "/" + image11, ProjectKind, "", ResourceNotFound}, // case 20
		{"/", project1 + "/" + image11, EnvKind, "", ResourceNotFound},
		{"/", project2 + "/" + image11, ImageKind, "", ResourceNotFound},
		{"/" + project1, project1 + "/" + image11, ImageKind, project1 + "/" + image11, nil},
		{"/" + project1 + "/" + image11, project1 + "/" + image11, ImageKind, project1 + "/" + image11, nil},
		{"/" + project2, project1 + "/" + image11, ImageKind, project1 + "/" + image11, nil},
		{"/" + project2 + "/" + image21, project1 + "/" + image11, ImageKind, project1 + "/" + image11, nil},

		// Image relative
		{"/" + project1, image11, ImageKind, project1 + "/" + image11, nil},
		{"/" + project2, image11, ImageKind, "", ResourceNotFound},

		// Resolving not existing resources
		{"/", "notExisting", ProjectKind, "", ResourceNotFound},
		{"/", "notExisting", EnvKind, "", ResourceNotFound}, // case 30
		{"/", "notExisting", ImageKind, "", ResourceNotFound},
	}

	for i, c := range cases {
		path := filepath.Join(fakeWorkspacePath, c.fromPath)
		err := os.Chdir(path)
		require.NoError(t, err, "should not error for chdir on case %d", i)
		r, err := resolveContextualResource(c.exprIn, c.kindIn)
		if c.errWanted == nil {
			assert.NoError(t, err, "should not error on case %d", i)
		} else {
			assert.ErrorIs(t, err, c.errWanted, "bad error for case %d", i)
		}

		if c.resNameWanted == "" {
			assert.Nil(t, r, "should not found a resource for case %d", i)
		} else {
			require.NotNil(t, r, "should found a resource for case %d", i)
			assert.Equal(t, c.resNameWanted, r.Name(), "Bad resource name for case %d", i)
		}
	}

}

func TestResolveResource(t *testing.T) {
	fakeWorkspacePath := initWorkspace(t)

	cases := []struct {
		fromPath, exprIn string
		kindIn Kind
		resNameWanted string
		errWanted error
	} {
		// --- Resolving not typed resource
		// Project
		{"/", project1, AllKind, "", InconsistentExpression}, // case 0
		{"/", project1, ProjectKind, project1, nil},
		{project1, project1, ProjectKind, project1, nil},
		{project2, project1, ProjectKind, project1, nil},
		{envDir + env1, project1, ProjectKind, project1, nil},
		{"..", project1, ProjectKind, "", settings.PathNotFound},

		// Env
		{"/", env1, EnvKind, env1, nil},
		{envDir + env1, env1, EnvKind, env1, nil},
		{envDir + env2, env1, EnvKind, env1, nil},
		{project1, env1, EnvKind, env1, nil},
		{"..", env1, EnvKind, "", settings.PathNotFound},

		// Image
		{"/", image11, ImageKind, "", ResourceNotFound}, // case 10
		{"/", project1 + "/" + image11, ImageKind, project1 + "/" + image11, nil},
		{project1, image11, ImageKind, project1 + "/" + image11, nil},
		{project2, image11, ImageKind, "", ResourceNotFound},
		{envDir + env1, image11, ImageKind, "", ResourceNotFound},
		{project1, image11, EnvKind, "", ResourceNotFound},
		{"..", image11, ImageKind, "", settings.PathNotFound},


		// Resolving typed resource
		// Project
		{"/", "p/" + project1, AllKind, project1, nil},
		{"/", "p/" + project1, ProjectKind, project1, nil},
		{project1, "p/" + project1, ProjectKind, project1, nil},
		{project2, "p/" + project1, ProjectKind, project1, nil}, // case 20

		// Env 
		{"/", "e/" + env1, AllKind, env1, nil},
		{"/", "e/" + env1, EnvKind, env1, nil},
		{envDir + env1, "e/" + env1, EnvKind, env1, nil},
		{envDir + env2, "e/" + env1, EnvKind, env1, nil},

		// Image 
		{"/", "i/" + project1 + "/" + image11, AllKind, project1 + "/" + image11, nil},
		{"/", "i/" + project1 + "/" + image11, ImageKind, project1 + "/" + image11, nil},
		{project1, "i/" + image11, ImageKind, project1 + "/" + image11, nil},
		{project1, "i/" + image21, ImageKind, "", ResourceNotFound},

		// Resolving not existing resources
		{"/", "notExisting", ProjectKind, "", ResourceNotFound},
		{"/", "/notExisting", ProjectKind, "", ResourceNotFound},
		{"/", "project/notExisting", ProjectKind, "", ResourceNotFound},
		{"/", "env/notExisting", EnvKind, "", ResourceNotFound},
		{"/", "image/notExisting/notExisting", ImageKind, "", ResourceNotFound},
		{project1, "notExisting", ProjectKind, "", ResourceNotFound},
		{project1, "/notExisting", ProjectKind, "", ResourceNotFound},
		{project1, "project/notExisting", ProjectKind, "", ResourceNotFound},
		{project1, "env/notExisting", EnvKind, "", ResourceNotFound},
		{project1, "image/notExisting/notExisting", ImageKind, "", ResourceNotFound},
		{"..", "notExisting", ProjectKind, "", settings.PathNotFound},
		{"..", "/notExisting", ProjectKind, "", settings.PathNotFound},
		{"..", "project/notExisting", ProjectKind, "", settings.PathNotFound},

		// Resolving bad kind resources
		{"/", "project/" + project1, EnvKind, "", InconsistentExpressionType},
		{"/", "env/" + env1, ProjectKind, "", InconsistentExpressionType},
		{"/", "image/" + project1 + "/" + image11, ProjectKind, "", InconsistentExpressionType},
		{project1, "project/" + project1, EnvKind, "", InconsistentExpressionType},
		{"/", "env/" + env1, ProjectKind, "", InconsistentExpressionType},
		{"/", "env/" + env1, ImageKind, "", InconsistentExpressionType},

	}

	for i, c := range cases {
		path := filepath.Join(fakeWorkspacePath, c.fromPath)
		err := os.Chdir(path)
		require.NoError(t, err, "should not error for chdir on case %d", i)
		r, err := resolveResource(c.exprIn, c.kindIn)
		if c.errWanted == nil {
			assert.NoError(t, err, "should not error on case %d", i)
		} else {
			assert.ErrorIs(t, err, c.errWanted, "bad error for case %d", i)
		}

		if c.resNameWanted == "" {
			assert.Nil(t, r, "should not found a resource for case %d", i)
		} else {
			require.NotNil(t, r, "should found a resource for case %d", i)
			assert.Equal(t, c.resNameWanted, r.Name(), "Bad resource name for case %d", i)
		}
	}
}

func TestSplitExpressions(t *testing.T) {
	cases := []struct {
                exprIn string
                strippedExprsWanted []string
		kindsWanted []Kind
                errWanted error
        } {
		{"p foo bar", []string{"foo", "bar"}, []Kind{ProjectKind}, nil}, // case 0
		{"project foo bar", []string{"foo", "bar"}, []Kind{ProjectKind}, nil},
		{"projects foo bar", []string{"foo", "bar"}, []Kind{ProjectKind}, nil},
		{"p p/foo project/bar", []string{"p/foo", "project/bar"}, []Kind{ProjectKind}, nil},
		{"p p/foo e/bar", []string{"p/foo", "e/bar"}, []Kind{ProjectKind}, nil},
		{"p,e p/foo e/bar", []string{"p/foo", "e/bar"}, []Kind{ProjectKind, EnvKind}, nil}, // case 5
		{"p,e,i foo bar", []string{"foo", "bar"}, []Kind{ProjectKind, EnvKind, ImageKind}, nil},

		{"foo bar baz", []string{"foo", "bar", "baz"}, nil, nil},
		{"a p/foo project/bar", []string{"a", "p/foo", "project/bar"}, nil, nil},
		{"p,a p/foo project/bar", nil, nil, UnknownKind},
		{"a,p p/foo project/bar", nil, nil, UnknownKind}, // case 10
	}

	for i, c := range cases {
		exprs, kinds, err := splitExpressions(c.exprIn)
		if c.errWanted == nil {
			assert.NoError(t, err, "should not error on case %d", i)
		} else {
			assert.ErrorIs(t, err, c.errWanted, "bad error for case %d", i)
		}

		assert.Equal(t, c.strippedExprsWanted, exprs, "bad stripped expressions for case %d", i)
		assert.Equal(t, c.kindsWanted, kinds, "bad kinds for case %d", i)
	}
}

func TestResolveExpression(t *testing.T) {
	fakeWorkspacePath := initWorkspace(t)

	cases := []struct {
		fromPath, exprIn string
		kindsIn []Kind
		resNamesWanted []string
		errWanted error
	} {
		// Projects resolution
		{"/", project1, []Kind{}, []string{}, InconsistentExpression}, // case 0
		{"/", project1, []Kind{AllKind}, []string{}, InconsistentExpression},
		{"/", "p/" + project1, []Kind{AllKind}, []string{project1}, nil},
		{"/", "p " + project1, []Kind{AllKind}, []string{project1}, nil},
		{"/", project1, []Kind{ProjectKind}, []string{project1}, nil},
		{"/", project1, []Kind{EnvKind}, []string{}, ResourceNotFound},
		{"/", project1 + " " + project2, []Kind{AllKind}, []string{}, InconsistentExpression},
		{"/", "p/" + project1 + " p/" + project2, []Kind{AllKind}, []string{project1, project2}, nil},
		{"/", "p " + project1 + " " + project2, []Kind{AllKind}, []string{project1, project2}, nil},
		{"/", project1 + " " + project2, []Kind{ProjectKind}, []string{project1, project2}, nil},

		// Envs resolution
		{"/", env1, []Kind{}, []string{}, InconsistentExpression}, // case 10
		{"/", env1, []Kind{AllKind}, []string{}, InconsistentExpression},
		{"/", "e/" + env1, []Kind{AllKind}, []string{env1}, nil},
		{"/", "e " + env1, []Kind{AllKind}, []string{env1}, nil},
		{"/", env1, []Kind{EnvKind}, []string{env1}, nil},
		{"/", env1 + " " + env2, []Kind{AllKind}, []string{}, InconsistentExpression},
		{"/", "e/" + env1 + " e/" + env2, []Kind{AllKind}, []string{env1, env2}, nil},
		{"/", "e " + env1 + " " + env2, []Kind{AllKind}, []string{env1, env2}, nil},
		{"/", env1 + " " + env2, []Kind{EnvKind}, []string{env1, env2}, nil},

		// Images resolution
		{"/", project2 + "/" + image21, []Kind{}, []string{}, InconsistentExpression},
		{"/", project2 + "/" + image21, []Kind{AllKind}, []string{}, InconsistentExpression}, // case 20
		{"/", "i/" + project2 + "/" + image21, []Kind{AllKind}, []string{project2 + "/" + image21}, nil},
		{"/", "i " + project2 + "/" + image21, []Kind{AllKind}, []string{project2 + "/" + image21}, nil},
		{"/", project2 + "/" + image21, []Kind{ImageKind}, []string{project2 + "/" + image21}, nil},
		{"/", project2 + "/" + image21 + " " + project1 + "/" + image12, []Kind{AllKind}, []string{}, InconsistentExpression},
		{"/", "i/" + project2 + "/" + image21 + " i/" + project1 + "/" + image12, []Kind{AllKind}, []string{project1 + "/" + image12, project2 + "/" + image21}, nil},
		{"/", "i " + project2 + "/" + image21 + " " + project1 + "/" + image12, []Kind{AllKind}, []string{project1 + "/" + image12, project2 + "/" + image21}, nil},
		{"/", project2 + "/" + image21 + " " + project1 + "/" + image12, []Kind{ImageKind}, []string{project1 + "/" + image12, project2 + "/" + image21}, nil},

		// Mixed resolution "AllKind"
		{"/", "p/" + project1 + " e/" + env2, []Kind{}, []string{project1, env2}, nil},
		{"/", "all p/" + project1 + " e/" + env2, []Kind{}, []string{project1, env2}, nil},
		{"/", "all " + project1 + " " + env2, []Kind{}, []string{}, InconsistentExpression}, // case 30
		{"/", "p/" + project1 + " e/" + env2, []Kind{AllKind}, []string{project1, env2}, nil},
		{"/", "all p/" + project1 + " e/" + env2, []Kind{AllKind}, []string{project1, env2}, nil},
		{"/", "all " + project1 + " " + env2, []Kind{AllKind}, []string{}, InconsistentExpression},
		{"/", "projects p/" + project1 + " e/" + env2, []Kind{AllKind}, []string{project1}, InconsistentExpressionType},
		{"/", "p " + project1 + " " + env2, []Kind{ProjectKind}, []string{project1}, ResourceNotFound},
		{"/", "projects,envs " + project1 + " " + env2, []Kind{AllKind}, []string{project1, env2}, nil},
		{"/", "projects,envs p/" + project1 + " e/" + env2, []Kind{AllKind}, []string{project1, env2}, nil},
		{"/", "p/" + project1 + " e/" + env2 + " i/" + project2 + "/" + image21, []Kind{AllKind}, []string{project1, env2, project2 + "/" + image21}, nil},
		{"/", "p,e,i p/" + project1 + " e/" + env2 + " i/" + project2 + "/" + image21, []Kind{AllKind}, []string{project1, env2, project2 + "/" + image21}, nil},
		{"/", "all p/" + project1 + " e/" + env2 + " i/" + project2 + "/" + image21, []Kind{AllKind}, []string{project1, env2, project2 + "/" + image21}, nil}, // case 40
		{"/", "p,e,i " + project1 + " " + env2 + " " + project2 + "/" + image21, []Kind{AllKind}, []string{project1, env2, project2 + "/" + image21}, nil},
		{"/", "p,e " + project1 + " " + env2 + " " + project2 + "/" + image21, []Kind{AllKind}, []string{project1, env2}, ResourceNotFound},
		{"/", "p,e " + project1 + " " + env2 + " i/" + project2 + "/" + image21, []Kind{AllKind}, []string{project1, env2}, InconsistentExpressionType},

		// Mixed resolution "Multiple kinds"
		{"/", "p/" + project1 + " e/" + env2, []Kind{ProjectKind, EnvKind}, []string{project1, env2}, nil},
		{"/", "projects p/" + project1 + " e/" + env2, []Kind{ProjectKind, EnvKind, ImageKind}, []string{project1}, InconsistentExpressionType},
		{"/", "p " + project1 + " " + env2, []Kind{ProjectKind, EnvKind}, []string{project1}, ResourceNotFound},
		{"/", "projects,envs " + project1 + " " + env2, []Kind{ProjectKind, EnvKind, ImageKind}, []string{project1, env2}, nil},
		{"/", "projects,envs p/" + project1 + " e/" + env2, []Kind{ProjectKind, EnvKind, ImageKind}, []string{project1, env2}, nil},
		{"/", "p/" + project1 + " e/" + env2 + " i/" + project2 + "/" + image21, []Kind{ProjectKind, EnvKind, ImageKind}, []string{project1, env2, project2 + "/" + image21}, nil},
		{"/", "p,e,i p/" + project1 + " e/" + env2 + " i/" + project2 + "/" + image21, []Kind{ProjectKind, EnvKind, ImageKind}, []string{project1, env2, project2 + "/" + image21}, nil}, // case 50
		{"/", "all p/" + project1 + " e/" + env2 + " i/" + project2 + "/" + image21, []Kind{ProjectKind, EnvKind, ImageKind}, []string{project1, env2, project2 + "/" + image21}, nil},
		{"/", "p,e,i " + project1 + " " + env2 + " " + project2 + "/" + image21, []Kind{ProjectKind, EnvKind, ImageKind}, []string{project1, env2, project2 + "/" + image21}, nil},
		{"/", "p,e " + project1 + " " + env2 + " " + project2 + "/" + image21, []Kind{ProjectKind, EnvKind, ImageKind}, []string{project1, env2}, ResourceNotFound},
		{"/", "p,e " + project1 + " " + env2 + " i/" + project2 + "/" + image21, []Kind{ProjectKind, EnvKind, ImageKind}, []string{project1, env2}, InconsistentExpressionType},

		// Mixed resolution and errors
		{"/", "p/notExist p/" + project1 + " e/" + env2, []Kind{ProjectKind, EnvKind}, []string{project1, env2}, ResourceNotFound},
		{"/", "p/" + project1 + " e/" + env2 + " p/notExist", []Kind{ProjectKind, EnvKind}, []string{project1, env2}, ResourceNotFound},
		{"/", "i/notExist p/" + project1 + " e/" + env2, []Kind{ProjectKind, EnvKind}, []string{project1, env2}, InconsistentExpressionType},
		{"/", "p/" + project1 + " e/" + env2 + " i/notExist", []Kind{ProjectKind, EnvKind}, []string{project1, env2}, InconsistentExpressionType},
		{"/", "p/" + project1 + " e/" + env2, []Kind{ProjectKind}, []string{project1}, InconsistentExpressionType},
		{"/", "p/" + project1 + " e/" + env2, []Kind{EnvKind}, []string{env2}, InconsistentExpressionType}, // case 60
	}

	for i, c := range cases {
		path := filepath.Join(fakeWorkspacePath, c.fromPath)
		err := os.Chdir(path)
		require.NoError(t, err, "should not error for chdir on case %d", i)
		resources, aggErr := ResolveExpression(c.exprIn, c.kindsIn...)
		if c.errWanted == nil {
			assert.Error(t, aggErr, "should return no aggregated error on case %d", i)
			assert.False(t, aggErr.GotError(), "should return no aggregated error on case %d", i)
			assert.Equal(t, "", aggErr.Error(), "should return no aggregated error on case %d", i)
		} else {
			//assert.ErrorIs(t, err, c.errWanted, "bad error for case %d", i)
			assert.True(t, aggErr.GotError(), "should return aggregated error on case %d", i)
			assert.True(t, aggErr.Is(c.errWanted), "bad error for case %d", i)
		}

		if len(c.resNamesWanted) == 0 {
			assert.Nil(t, resources, "should not found a resource for case %d", i)
		} else {
			require.NotNil(t, resources, "should found some resources for case %d", i)
			assert.Len(t, resources, len(c.resNamesWanted), "bad resources count returned for case %d", i)
			for _, res := range resources {
				require.NotNil(t, res, "should found not nil resource for case %d", i)
				assert.Contains(t, c.resNamesWanted, res.Name(), "Bad resource name [%s] for case %d", res.Name(), i)
			}
		}
	}
}
