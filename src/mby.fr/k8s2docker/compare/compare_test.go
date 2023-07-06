package compare

import (
	"testing"
	k8sv1 "k8s.io/api/core/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	pod1a_descriptor := `
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
	
	pod1b_descriptor := `
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
)

func loadPod(t *testing.T, s string) (pod k8sv1.Pod) {
	err = return json.Unmarshal(s, &pod)
	require.NoError(t, err, fmt.Sprintf("Unable to load pod %s", s))
	require.NotNil(t, pod, fmt.Sprintf("Unable to load pod %s", s))
	return
}

func TestComparePods_Identicals(t *testing.T) {
	poda := loadPod(pod1a_descriptor)
	podb := loadPod(pod1a_descriptor)

	podDiff := comparePods(poda, podb)
	require.NotNil(podDiff)
	assert.False(t, podDiff.DoDiffer())
	assert.False(t, podDiff.IsRestartNeeded())
	assert.Len(t, 0, podDiff.Diffs())
	assert.Len(t, 0, podDiff.UpdatableDiffs())
}

func TestComparePods_Differents(t *testing.T) {
	pod1a := loadPod(pod1a_descriptor)
	pod1b := loadPod(pod1b_descriptor)

	podDiff := comparePods(pod1a, pod1b)
	require.NotNil(podDiff)
	assert.True(t, podDiff.DoDiffer())
	assert.True(t, podDiff.IsRestartNeeded())
	assert.Len(t, 1, podDiff.Diffs())
	assert.Len(t, 0, podDiff.UpdatableDiffs())
}

