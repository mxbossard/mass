package repo

import (
	"gopkg.in/yaml.v3"
	_ "io"
	_ "net/http"
	_ "net/http/httptest"
	"os"
	"testing"

	"mby.fr/utils/serializ"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testDbPath = "./testMyDb"

	expectedNsYaml1 = `
apiVersion: v1
kind: Namespace
metadata:
  name: ns1
`
	expectedNsYaml2 = `
apiVersion: v1
kind: Namespace
metadata:
  name: ns2
`
	expectedNsYaml3 = `
apiVersion: v1
kind: Namespace
metadata:
  name: ns3
`

	expectedPodYaml1 = `
apiVersion: v1
kind: Pod
metadata:
  name: pod1
`
	expectedPodYaml2 = `
apiVersion: v1
kind: Pod
metadata:
  name: pod2
`
	expectedPodYaml3 = `
apiVersion: v1
kind: Pod
metadata:
  name: pod3
`

	expectedPodYaml1InNs1 = expectedPodYaml1 + `
  namespace: ns1
`
	expectedPodYaml2InNs1 = expectedPodYaml2 + `
  namespace: ns1
`
	expectedPodYaml3InNs2 = expectedPodYaml3 + `
  namespace: ns2
`

	expectedServiceYaml1 = `
apiVersion: v1
kind: Service
metadata:
  name: svc1
`
	expectedServiceYaml2 = `
apiVersion: v1
kind: Service
metadata:
  name: svc2
`
	expectedServiceYaml3 = `
apiVersion: v1
kind: Service
metadata:
  name: svc3
`
)

func initTestDb() {
	initDb(testDbPath)
}

func clearTestDb() {
	os.RemoveAll(testDbPath)
}

func mapYaml(t *testing.T, yamlIn string) map[string]any {
	var tree map[string]any
	err := yaml.Unmarshal([]byte(yamlIn), &tree)
	require.NoErrorf(t, err, "Unable to map yaml: [%s]", yamlIn)
	return tree
}

func yamlToJson(t *testing.T, yamlIn string) string {
	json, err := serializ.YamlToJsonString(yamlIn)
	require.NoErrorf(t, err, "Unable to convert to json: [%s]", yamlIn)
	return json
}

func storeYamlResource(t *testing.T, namespace, yamlIn string) {
	tree := mapYaml(t, yamlIn)
	_, err := storeResource(namespace, tree)
	require.NoErrorf(t, err, "Unable to store res in namespace %s : from yaml: [%s]", namespace, yamlIn)
}

func TestConsolidateMetadata(t *testing.T) {
	var ns, k, n string
	ns, k, n = consolidateMetadata("", "", "", "")
	assert.Equal(t, "", ns)
	assert.Equal(t, "", k)
	assert.Equal(t, "", n)

	ns, k, n = consolidateMetadata("foo", "", "", "")
	assert.Equal(t, "foo", ns)
	assert.Equal(t, "", k)
	assert.Equal(t, "", n)

	ns, k, n = consolidateMetadata("", "bar", "", "")
	assert.Equal(t, "", ns)
	assert.Equal(t, "bar", k)
	assert.Equal(t, "", n)

	ns, k, n = consolidateMetadata("", "", "baz", "")
	assert.Equal(t, "", ns)
	assert.Equal(t, "", k)
	assert.Equal(t, "baz", n)

	ns, k, n = consolidateMetadata("foo", "bar", "baz", "")
	assert.Equal(t, "foo", ns)
	assert.Equal(t, "bar", k)
	assert.Equal(t, "baz", n)

	ns, k, n = consolidateMetadata("", "", "", yamlToJson(t, expectedPodYaml1))
	assert.Equal(t, "", ns)
	assert.Equal(t, "Pod", k)
	assert.Equal(t, "pod1", n)

	ns, k, n = consolidateMetadata("", "", "", yamlToJson(t, expectedPodYaml1))
	assert.Equal(t, "", ns)
	assert.Equal(t, "Pod", k)
	assert.Equal(t, "pod1", n)

	ns, k, n = consolidateMetadata("", "", "", yamlToJson(t, expectedPodYaml1InNs1))
	assert.Equal(t, "ns1", ns)
	assert.Equal(t, "Pod", k)
	assert.Equal(t, "pod1", n)

	ns, k, n = consolidateMetadata("foo", "", "", yamlToJson(t, expectedPodYaml1))
	assert.Equal(t, "foo", ns)
	assert.Equal(t, "Pod", k)
	assert.Equal(t, "pod1", n)

	ns, k, n = consolidateMetadata("foo", "", "", yamlToJson(t, expectedPodYaml1InNs1))
	assert.Equal(t, "ns1", ns)
	assert.Equal(t, "Pod", k)
	assert.Equal(t, "pod1", n)

	ns, k, n = consolidateMetadata("foo", "bar", "", yamlToJson(t, expectedPodYaml1InNs1))
	assert.Equal(t, "ns1", ns)
	assert.Equal(t, "Pod", k)
	assert.Equal(t, "pod1", n)

	ns, k, n = consolidateMetadata("foo", "bar", "baz", yamlToJson(t, expectedPodYaml1InNs1))
	assert.Equal(t, "ns1", ns)
	assert.Equal(t, "Pod", k)
	assert.Equal(t, "pod1", n)
}

func TestCompleteJsonInput(t *testing.T) {
	// TODO
}

func TestDevelopNamespaceNames(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	res, err := developNamespaceNames("")
	require.NoError(t, err)
	assert.Len(t, res, 0)

	res, err = developNamespaceNames("foo")
	require.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Contains(t, res, "foo")

	res, err = developNamespaceNames("all")
	require.NoError(t, err)
	assert.Len(t, res, 0)

	storeYamlResource(t, "", expectedNsYaml1)
	require.NoError(t, err)
	storeYamlResource(t, "", expectedNsYaml2)
	require.NoError(t, err)
	res, err = developNamespaceNames("all")
	require.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Contains(t, res, "ns1")
	assert.Contains(t, res, "ns2")
}

func TestListResourcesAsMap_Namespaces(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	var res []map[string]any
	var err error
	res, err = listResourcesAsMap("", "", "")
	require.NoError(t, err)
	assert.Len(t, res, 0)

	// Add some Namespaces in repo
	storeYamlResource(t, "", expectedNsYaml1)
	require.NoError(t, err)
	res, err = listResourcesAsMap("", "", "")
	require.NoError(t, err)
	assert.Len(t, res, 1)

	storeYamlResource(t, "", expectedNsYaml2)
	storeYamlResource(t, "", expectedNsYaml3)
	res, err = listResourcesAsMap("", "", "")
	require.NoError(t, err)
	assert.Len(t, res, 3)
}

func TestListResourcesAsMap_Namespace_AllKinds(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	var res []map[string]any
	var err error

	// Add some Namespaces in repo
	//storeYamlResource(t, "", expectedNsYaml1)
	//storeYamlResource(t, "", expectedNsYaml2)

	// List all resources in ns1
	res, err = listResourcesAsMap("ns1", "", "")
	require.NoError(t, err)
	assert.Len(t, res, 0)

	// Add some resources in repo
	storeYamlResource(t, "ns1", expectedPodYaml1)
	storeYamlResource(t, "ns1", expectedPodYaml2)
	storeYamlResource(t, "ns2", expectedPodYaml3)
	storeYamlResource(t, "ns1", expectedServiceYaml1)
	storeYamlResource(t, "ns1", expectedServiceYaml2)
	storeYamlResource(t, "ns2", expectedServiceYaml3)

	res, err = listResourcesAsMap("ns1", "", "")
	require.NoError(t, err)
	assert.Len(t, res, 4)
	assert.Contains(t, res, mapYaml(t, expectedPodYaml1))
	assert.Contains(t, res, mapYaml(t, expectedPodYaml2))
	assert.Contains(t, res, mapYaml(t, expectedServiceYaml1))
	assert.Contains(t, res, mapYaml(t, expectedServiceYaml2))

	res, err = listResourcesAsMap("ns2", "", "")
	require.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Contains(t, res, mapYaml(t, expectedPodYaml3))
	assert.Contains(t, res, mapYaml(t, expectedServiceYaml3))
}

func TestListResourcesAsMap_AllNamespaces_AllKinds(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	var res []map[string]any
	var err error

	res, err = listResourcesAsMap("all", "", "")
	require.NoError(t, err)
	assert.Len(t, res, 0)

	// Add some resources in repo
	storeYamlResource(t, "ns1", expectedPodYaml1)
	storeYamlResource(t, "ns1", expectedPodYaml2)
	storeYamlResource(t, "ns2", expectedPodYaml3)
	storeYamlResource(t, "ns1", expectedServiceYaml1)
	storeYamlResource(t, "ns1", expectedServiceYaml2)
	storeYamlResource(t, "ns2", expectedServiceYaml3)

	res, err = listResourcesAsMap("all", "", "")
	require.NoError(t, err)
	assert.Len(t, res, 6)
	assert.Contains(t, res, mapYaml(t, expectedPodYaml1))
	assert.Contains(t, res, mapYaml(t, expectedPodYaml2))
	assert.Contains(t, res, mapYaml(t, expectedServiceYaml1))
	assert.Contains(t, res, mapYaml(t, expectedServiceYaml2))
	assert.Contains(t, res, mapYaml(t, expectedPodYaml3))
	assert.Contains(t, res, mapYaml(t, expectedServiceYaml3))
}

func TestListResourcesAsMap_Namespace_Kind(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	var res []map[string]any
	var err error

	res, err = listResourcesAsMap("ns1", "Pod", "")
	require.NoError(t, err)
	assert.Len(t, res, 0)

	// Add some resources in repo
	storeYamlResource(t, "ns1", expectedPodYaml1)
	storeYamlResource(t, "ns1", expectedPodYaml2)
	storeYamlResource(t, "ns2", expectedPodYaml3)
	storeYamlResource(t, "ns1", expectedServiceYaml1)
	storeYamlResource(t, "ns1", expectedServiceYaml2)
	storeYamlResource(t, "ns2", expectedServiceYaml3)

	res, err = listResourcesAsMap("ns1", "Pod", "")
	require.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Contains(t, res, mapYaml(t, expectedPodYaml1))
	assert.Contains(t, res, mapYaml(t, expectedPodYaml2))
}

func TestListResourcesAsMap_AllNamespaces_Kind(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	var res []map[string]any
	var err error

	res, err = listResourcesAsMap("all", "Pod", "")
	require.NoError(t, err)
	assert.Len(t, res, 0)

	// Add some resources in repo
	storeYamlResource(t, "ns1", expectedPodYaml1)
	storeYamlResource(t, "ns1", expectedPodYaml2)
	storeYamlResource(t, "ns2", expectedPodYaml3)
	storeYamlResource(t, "ns1", expectedServiceYaml1)
	storeYamlResource(t, "ns1", expectedServiceYaml2)
	storeYamlResource(t, "ns2", expectedServiceYaml3)

	res, err = listResourcesAsMap("all", "Pod", "")
	require.NoError(t, err)
	assert.Len(t, res, 3)
	assert.Contains(t, res, mapYaml(t, expectedPodYaml1))
	assert.Contains(t, res, mapYaml(t, expectedPodYaml2))
	assert.Contains(t, res, mapYaml(t, expectedPodYaml3))
}

func TestListResourcesAsMap_Namespace_Kind_Name(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	var res []map[string]any
	var err error

	res, err = listResourcesAsMap("all", "Pod", "pod1")
	require.NoError(t, err)
	assert.Len(t, res, 0)

	// Add some resources in repo
	storeYamlResource(t, "ns1", expectedPodYaml1)
	storeYamlResource(t, "ns1", expectedPodYaml2)
	storeYamlResource(t, "ns2", expectedPodYaml3)
	storeYamlResource(t, "ns1", expectedServiceYaml1)
	storeYamlResource(t, "ns1", expectedServiceYaml2)
	storeYamlResource(t, "ns2", expectedServiceYaml3)

	res, err = listResourcesAsMap("ns1", "Pod", "pod1")
	require.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Contains(t, res, mapYaml(t, expectedPodYaml1))

	res, err = listResourcesAsMap("ns1", "Pod", "pod2")
	require.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Contains(t, res, mapYaml(t, expectedPodYaml2))
}

func TestListResourcesAsMap_AllNamespaces_Kind_Name(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	var res []map[string]any
	var err error

	res, err = listResourcesAsMap("all", "Pod", "pod1")
	require.NoError(t, err)
	assert.Len(t, res, 0)

	// Add some resources in repo
	storeYamlResource(t, "ns1", expectedPodYaml1)
	storeYamlResource(t, "ns1", expectedPodYaml2)
	storeYamlResource(t, "ns2", expectedPodYaml1)
	storeYamlResource(t, "ns1", expectedServiceYaml1)
	storeYamlResource(t, "ns1", expectedServiceYaml2)
	storeYamlResource(t, "ns2", expectedServiceYaml3)

	res, err = listResourcesAsMap("all", "Pod", "pod1")
	require.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Contains(t, res, mapYaml(t, expectedPodYaml1))

	res, err = listResourcesAsMap("all", "Pod", "pod2")
	require.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Contains(t, res, mapYaml(t, expectedPodYaml2))
}

func TestGet_Namespace(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	// List all namespaces
	out, err := Get("", "", "", "")
	require.NoError(t, err)
	assert.Empty(t, out)

	storeYamlResource(t, "", expectedNsYaml1)
	storeYamlResource(t, "", expectedNsYaml2)
	storeYamlResource(t, "", expectedNsYaml3)

	// List all namespaces
	out, err = Get("", "", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml1))
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml2))
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml3))

	// List all resources of ns1
	out, err = Get("ns1", "", "", "")
	require.NoError(t, err)
	assert.Empty(t, out)

	storeYamlResource(t, "ns1", expectedPodYaml1)
	storeYamlResource(t, "ns1", expectedPodYaml2)
	storeYamlResource(t, "ns2", expectedPodYaml3)
	storeYamlResource(t, "ns1", expectedServiceYaml1)
	storeYamlResource(t, "ns1", expectedServiceYaml2)
	storeYamlResource(t, "ns2", expectedServiceYaml3)

	// List all resources of ns1
	out, err = Get("ns1", "", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml1))
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml2))
	assert.NotContains(t, out, yamlToJson(t, expectedPodYaml3))
	assert.Contains(t, out, yamlToJson(t, expectedServiceYaml1))
	assert.Contains(t, out, yamlToJson(t, expectedServiceYaml2))
	assert.NotContains(t, out, yamlToJson(t, expectedServiceYaml3))

	// List all resources of all namespaces
	out, err = Get("all", "", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml1))
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml2))
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml3))
	assert.Contains(t, out, yamlToJson(t, expectedServiceYaml1))
	assert.Contains(t, out, yamlToJson(t, expectedServiceYaml2))
	assert.Contains(t, out, yamlToJson(t, expectedServiceYaml3))
}

func TestGet_Kind(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	// List all Pods of ns1 namespace
	out, err := Get("ns1", "Pod", "", "")
	require.NoError(t, err)
	assert.Empty(t, out)

	storeYamlResource(t, "ns1", expectedPodYaml1)
	storeYamlResource(t, "ns1", expectedPodYaml2)
	storeYamlResource(t, "ns2", expectedPodYaml3)
	storeYamlResource(t, "ns1", expectedServiceYaml1)
	storeYamlResource(t, "ns1", expectedServiceYaml2)
	storeYamlResource(t, "ns2", expectedServiceYaml3)

	out, err = Get("ns1", "Pod", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml1))
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml2))
	assert.NotContains(t, out, yamlToJson(t, expectedPodYaml3))

	out, err = Get("all", "Pod", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml1))
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml2))
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml3))

	// FIXME: what to do ?
	out, err = Get("", "Pod", "", "")
	_ = out
	_ = err
}

func TestGet_Name(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	// List all Pods of ns1 namespace
	out, err := Get("ns1", "Pod", "pod1", "")
	require.NoError(t, err)
	assert.Empty(t, out)

	storeYamlResource(t, "ns1", expectedPodYaml1)
	storeYamlResource(t, "ns1", expectedPodYaml2)
	storeYamlResource(t, "ns2", expectedPodYaml3)
	storeYamlResource(t, "ns1", expectedServiceYaml1)
	storeYamlResource(t, "ns1", expectedServiceYaml2)
	storeYamlResource(t, "ns2", expectedServiceYaml3)

	out, err = Get("ns1", "Pod", "pod1", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml1))
	assert.NotContains(t, out, yamlToJson(t, expectedPodYaml2))
	assert.NotContains(t, out, yamlToJson(t, expectedPodYaml3))

	out, err = Get("ns2", "Pod", "pod1", "")
	require.NoError(t, err)
	assert.Empty(t, out)

	out, err = Get("ns2", "Pod", "pod3", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.NotContains(t, out, yamlToJson(t, expectedPodYaml1))
	assert.NotContains(t, out, yamlToJson(t, expectedPodYaml2))
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml3))

	// FIXME: what to do ?
	out, err = Get("", "Pod", "pod1", "")
	_ = out
	_ = err

	// FIXME: what to do ?
	out, err = Get("", "", "pod1", "")
	_ = out
	_ = err
}

func TestPost_Namespace(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	out, err := Get("", "", "", "")
	require.NoError(t, err)
	assert.Empty(t, out)

	out, err = Post("ns1", "", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml1))
	assert.NotContains(t, out, yamlToJson(t, expectedNsYaml2))

	out, err = Get("", "", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml1))
	assert.NotContains(t, out, yamlToJson(t, expectedNsYaml2))

	out, err = Post("ns1", "", "", "")
	require.Error(t, err, "Overwriting should fail !")
	assert.Empty(t, out)

	out, err = Post("ns2", "", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml2))
	assert.NotContains(t, out, yamlToJson(t, expectedNsYaml1))

	out, err = Get("", "", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml1))
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml2))
}

func TestPost_Pod(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	out, err := Get("ns1", "Pod", "", "")
	require.NoError(t, err)
	assert.Empty(t, out)

	// Cannot post a pod without a name nor a json content
	out, err = Post("ns1", "Pod", "", "")
	require.Error(t, err, "Posting noname resource should fail !")
	assert.Empty(t, out)

	out, err = Get("ns1", "Pod", "", "")
	require.NoError(t, err)
	assert.Empty(t, out)

	out, err = Post("ns1", "Pod", "pod2", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml2InNs1))

	out, err = Get("ns1", "Pod", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml2InNs1))

	// Should return error because cannot overwrite nor update resource
	out, err = Post("ns1", "Pod", "pod2", "")
	require.Error(t, err, "Overwriting should fail !")
	assert.Empty(t, out)

	out, err = Post("ns1", "", "", yamlToJson(t, expectedPodYaml1))
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml1InNs1))

	out, err = Get("ns1", "Pod", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml1InNs1))

	// Should error because no namespace supplied nether arg nor in yaml
	out, err = Post("", "", "", yamlToJson(t, expectedPodYaml1))
	require.Error(t, err, "Posting not namespaced resource should fail !")
	assert.Empty(t, out)

	out, err = Post("", "", "", yamlToJson(t, expectedPodYaml3InNs2))
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml2))
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml3InNs2))

	out, err = Get("ns2", "Pod", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml3InNs2))
}

func TestPut_Namespace(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	out, err := Get("", "", "", "")
	require.NoError(t, err)
	assert.Empty(t, out)

	out, err = Put("ns1", "", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml1))
	assert.NotContains(t, out, yamlToJson(t, expectedNsYaml2))

	out, err = Get("", "", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml1))
	assert.NotContains(t, out, yamlToJson(t, expectedNsYaml2))

	// Overwrite ok with put
	out, err = Put("ns1", "", "", "")
	require.NoError(t, err, "Overwriting should succeed !")
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml1))
	assert.NotContains(t, out, yamlToJson(t, expectedNsYaml2))

	out, err = Put("ns2", "", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml2))
	assert.NotContains(t, out, yamlToJson(t, expectedNsYaml1))

	out, err = Get("", "", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml1))
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml2))
}

func TestPut_Pod(t *testing.T) {
	initTestDb()
	defer clearTestDb()

	out, err := Get("ns1", "Pod", "", "")
	require.NoError(t, err)
	assert.Empty(t, out)

	// Cannot put a pod without a name nor a json content
	out, err = Put("ns1", "Pod", "", "")
	require.Error(t, err, "Posting noname resource should fail !")
	assert.Empty(t, out)

	out, err = Get("ns1", "Pod", "", "")
	require.NoError(t, err)
	assert.Empty(t, out)

	out, err = Put("ns1", "Pod", "pod2", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml2InNs1))

	out, err = Get("ns1", "Pod", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml2InNs1))

	// Overwrite ok with put
	out, err = Put("ns1", "Pod", "pod2", "")
	require.NoError(t, err, "Overwriting should succeed !")
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml2InNs1))

	out, err = Put("ns1", "", "", yamlToJson(t, expectedPodYaml1))
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml1InNs1))

	out, err = Get("ns1", "Pod", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml1InNs1))

	// Should error because no namespace supplied nether arg nor in yaml
	out, err = Put("", "", "", yamlToJson(t, expectedPodYaml1))
	require.Error(t, err, "Puting not namespaced resource should fail !")
	assert.Empty(t, out)

	out, err = Put("", "", "", yamlToJson(t, expectedPodYaml3InNs2))
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedNsYaml2))
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml3InNs2))

	out, err = Get("ns2", "Pod", "", "")
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, yamlToJson(t, expectedPodYaml3InNs2))
}
