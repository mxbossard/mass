package compare

import (
	"fmt"
	//"gopkg.in/yaml.v3"
	//"encoding/json"
	k8sv1 "k8s.io/api/core/v1"
	//"sigs.k8s.io/yaml"
	//"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mby.fr/k8s2docker/descriptor"
)

var (
	empty_descriptor = `
apiVersion: v1
kind: Pod
`

	pod1a_descriptor = `
apiVersion: v1
kind: Pod
metadata:
  name: foo
spec:
  containers:
  - name: nging
    image: nginx:1.14.2
    ports:
    - containerPort: 80
`

	pod1b_descriptor = `
apiVersion: v1
kind: Pod
metadata:
  name: foo
spec:
  containers:
  - name: nging
    image: nginx:1.15.3
    ports:
    - containerPort: 80
`

	pod1c_descriptor = `
apiVersion: v1
kind: Pod
metadata:
  name: foo
  labels:
    foo: bar
spec:
  containers:
  - name: nging
    image: nginx:1.15.3
    ports:
    - containerPort: 80
  restartPolicy: Always
`

	pod1d_descriptor = `
apiVersion: v1
kind: Pod
metadata:
  name: foo
  labels:
    foo: bar
spec:
  containers:
  - name: baz
    image: nginx:1.15.3
    ports:
    - containerPort: 81
  - name: nging
    image: nginx:1.15.3
    ports:
    - containerPort: 80
  restartPolicy: Never
`
)

func loadPod(t *testing.T, s string) (pod k8sv1.Pod) {
	pod, err := descriptor.LoadPod([]byte(s))
	require.NoError(t, err, fmt.Sprintf("Unable to load pod %s", s))
	require.NotNil(t, pod, fmt.Sprintf("Unable to load pod %s", s))
	//fmt.Printf("Loaded pod: %v\n\n", pod)
	return
}

func TestAnyEquals(t *testing.T) {
	ok := anyEquals(nil, nil)
	assert.True(t, ok)

	ok = anyEquals("nil", nil)
	assert.False(t, ok)

	ok = anyEquals(nil, "nil")
	assert.False(t, ok)

	ok = anyEquals("nil", "nil")
	assert.True(t, ok)

	pod1a1 := loadPod(t, pod1a_descriptor)
	pod1a2 := loadPod(t, pod1a_descriptor)
	pod1b1 := loadPod(t, pod1b_descriptor)
	pod1c1 := loadPod(t, pod1c_descriptor)

	ok = anyEquals(pod1a1, pod1a2)
	assert.True(t, ok)

	ok = anyEquals(pod1a1, pod1b1)
	assert.False(t, ok)

	ok = anyEquals(pod1a1, pod1c1)
	assert.False(t, ok)

	ok = anyEquals(pod1b1, pod1c1)
	assert.False(t, ok)
}

func TestAppendDiff(t *testing.T) {
	var d []differ
	ok := appendDiff(&d, "foo", "", "")
	assert.False(t, ok)
	assert.Len(t, d, 0)

	d = nil
	ok = appendDiff(&d, "foo", "abc", "")
	assert.True(t, ok)
	require.Len(t, d, 1)
	assert.Equal(t, "foo", d[0].Path())
	assert.Equal(t, "abc", d[0].Left())
	assert.Equal(t, "", d[0].Right())

	d = nil
	ok = appendDiff(&d, "bar", "", "bcd")
	assert.True(t, ok)
	require.Len(t, d, 1)
	assert.Equal(t, "bar", d[0].Path())
	assert.Equal(t, "", d[0].Left())
	assert.Equal(t, "bcd", d[0].Right())

	d = nil
	ok = appendDiff(&d, "baz", "cde", "cde")
	assert.False(t, ok)
	require.Len(t, d, 0)

	d = nil
	ok = appendDiff(&d, "baz", "cde", "def")
	assert.True(t, ok)
	require.Len(t, d, 1)
	assert.Equal(t, "baz", d[0].Path())
	assert.Equal(t, "cde", d[0].Left())
	assert.Equal(t, "def", d[0].Right())
}

func TestComparePods_Empty(t *testing.T) {
	empty1 := loadPod(t, empty_descriptor)
	empty2 := loadPod(t, empty_descriptor)
	poda := loadPod(t, pod1a_descriptor)
	podb := loadPod(t, pod1b_descriptor)

	podDiff := ComparePods(empty1, empty2)
	require.NotNil(t, podDiff)
	assert.False(t, podDiff.DoDiffer())
	assert.Len(t, podDiff.Diffs(), 0)
	assert.Len(t, podDiff.updatableDiffs(), 0)
	assert.True(t, podDiff.IsUpdatable())

	podDiff = ComparePods(empty1, podb)
	require.NotNil(t, podDiff)
	assert.True(t, podDiff.DoDiffer())
	assert.Len(t, podDiff.Diffs(), 2)
	//assert.Len(t, podDiff.updatableDiffs(), 0)
	assert.False(t, podDiff.IsUpdatable())

	podDiff = ComparePods(poda, empty2)
	require.NotNil(t, podDiff)
	assert.True(t, podDiff.DoDiffer())
	assert.Len(t, podDiff.Diffs(), 2)
	//assert.Len(t, podDiff.updatableDiffs(), 0)
	assert.False(t, podDiff.IsUpdatable())
}

func TestComparePods_Identicals(t *testing.T) {
	poda := loadPod(t, pod1a_descriptor)
	podb := loadPod(t, pod1a_descriptor)

	podDiff := ComparePods(poda, podb)
	require.NotNil(t, podDiff)
	assert.False(t, podDiff.DoDiffer())
	assert.Len(t, podDiff.Diffs(), 0)
	assert.Len(t, podDiff.updatableDiffs(), 0)
	assert.True(t, podDiff.IsUpdatable())
}

func TestComparePods_Restartable(t *testing.T) {
	pod1a := loadPod(t, pod1a_descriptor)
	pod1b := loadPod(t, pod1b_descriptor)

	podDiff := ComparePods(pod1a, pod1b)
	require.NotNil(t, podDiff)
	assert.True(t, podDiff.DoDiffer())
	require.Len(t, podDiff.Diffs(), 1)
	assert.Len(t, podDiff.updatableDiffs(), 0)
	assert.False(t, podDiff.IsUpdatable())
	assert.Contains(t, podDiff.DiffPathes(), "pod.spec.containers.image")
}

func TestComparePods_Updatable(t *testing.T) {
	pod1b := loadPod(t, pod1b_descriptor)
	pod1c := loadPod(t, pod1c_descriptor)
	//log.Printf("pod1b: %v", pod1b.Spec.RestartPolicy)
	//log.Printf("pod1c: %v", pod1c.Spec.RestartPolicy)

	podDiff := ComparePods(pod1b, pod1c)
	require.NotNil(t, podDiff)
	assert.True(t, podDiff.DoDiffer())
	require.Len(t, podDiff.Diffs(), 2)
	assert.Len(t, podDiff.updatableDiffs(), 2)
	assert.True(t, podDiff.IsUpdatable())
	assert.Contains(t, podDiff.DiffPathes(), "pod.metadata.labels")
	assert.Contains(t, podDiff.DiffPathes(), "pod.spec.restartPolicy")
}

func TestComparePods_AddContainer(t *testing.T) {
	podA := loadPod(t, pod1c_descriptor)
	podB := loadPod(t, pod1d_descriptor)

	podDiff := ComparePods(podA, podB)
	require.NotNil(t, podDiff)
	assert.True(t, podDiff.DoDiffer())
	require.Len(t, podDiff.Diffs(), 2)
	assert.Len(t, podDiff.updatableDiffs(), 2)
	assert.True(t, podDiff.IsUpdatable())
	assert.Contains(t, podDiff.DiffPathes(), "pod.spec.containers")
	assert.Contains(t, podDiff.DiffPathes(), "pod.spec.restartPolicy")
	assert.Nil(t, podDiff.Diffs()[0].Left())
	assert.NotNil(t, podDiff.Diffs()[0].Right())
}

func TestComparePods_RemoveContainer(t *testing.T) {
	podA := loadPod(t, pod1c_descriptor)
	podB := loadPod(t, pod1d_descriptor)

	podDiff := ComparePods(podB, podA)
	require.NotNil(t, podDiff)
	assert.True(t, podDiff.DoDiffer())
	require.Len(t, podDiff.Diffs(), 2)
	assert.Len(t, podDiff.updatableDiffs(), 2)
	assert.True(t, podDiff.IsUpdatable())
	assert.Contains(t, podDiff.DiffPathes(), "pod.spec.containers")
	assert.Contains(t, podDiff.DiffPathes(), "pod.spec.restartPolicy")
	assert.Nil(t, podDiff.Diffs()[0].Right())
	assert.NotNil(t, podDiff.Diffs()[0].Left())
}
