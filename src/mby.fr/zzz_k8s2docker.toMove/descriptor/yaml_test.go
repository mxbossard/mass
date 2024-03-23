package descriptor

import (
	//"fmt"
	"io"
	"log"
	"strings"
	"testing"

	k8sv1 "k8s.io/api/core/v1"

	//"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	pod1a_descriptor = `
apiVersion: v1
kind: Pod
metadata:
  name: foo
spec:
  containers:
  - name: nginx
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
  - name: nginx
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
  - name: nginx
    image: nginx:1.15.3
    ports:
    - containerPort: 80
  restartPolicy: always
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
  - name: nginx
    image: nginx:1.15.3
    ports:
    - containerPort: 80
  restartPolicy: never
`
)

func TestValidator(t *testing.T) {
	v, err := validator.New(nil, validator.Opts{Strict: true})
	if err != nil {
		log.Fatalf("failed initializing validator: %s", err)
	}

	rc := io.NopCloser(strings.NewReader(pod1d_descriptor))
	for i, res := range v.Validate("foo", rc) { // A file might contain multiple resources
		// File starts with ---, the parser assumes a first empty resource
		if res.Status == validator.Invalid {
			log.Fatalf("resource %d is not valid: %s", i, res.Err)
		}
		if res.Status == validator.Error {
			log.Fatalf("error while processing resource %d: %s", i, res.Err)
		}
		//log.Printf("validation errors: %v\n", res.ValidationErrors)
		//log.Fatalf("Successfuly validated res: %v\n", string(res.Resource.Bytes))
	}
}

func TestLoadPod(t *testing.T) {
	pod, err := LoadPod([]byte(pod1a_descriptor))
	require.NoError(t, err, "Unable to load pod")
	assert.NotNil(t, pod, "Unable to load pod")
	assert.Equal(t, "foo", pod.ObjectMeta.Name)
	assert.Len(t, pod.Spec.Volumes, 0)
	assert.Len(t, pod.Spec.InitContainers, 0)
	require.Len(t, pod.Spec.Containers, 1)
	assert.Equal(t, "nginx", pod.Spec.Containers[0].Name)
	assert.Equal(t, "nginx:1.14.2", pod.Spec.Containers[0].Image)
	assert.Equal(t, k8sv1.RestartPolicy(""), pod.Spec.RestartPolicy)
}

func TestLoadPodWithDefaults_1a(t *testing.T) {
	t.Skip()
	pod, err := LoadPodWithDefaults([]byte(pod1a_descriptor))
	require.NoError(t, err, "Unable to load pod")
	assert.NotNil(t, pod, "Unable to load pod")
	assert.Equal(t, "foo", pod.ObjectMeta.Name)
	assert.Len(t, pod.Spec.Volumes, 0)
	assert.Len(t, pod.Spec.InitContainers, 0)
	require.Len(t, pod.Spec.Containers, 1)
	assert.Equal(t, "nginx", pod.Spec.Containers[0].Name)
	assert.Equal(t, "nginx:1.14.2", pod.Spec.Containers[0].Image)
	assert.Equal(t, k8sv1.RestartPolicyNever, pod.Spec.RestartPolicy)
}

func TestLoadPodWithDefaults_1b(t *testing.T) {
	t.Skip()
	pod, err := LoadPodWithDefaults([]byte(pod1b_descriptor))
	require.NoError(t, err, "Unable to load pod")
	assert.NotNil(t, pod, "Unable to load pod")
	assert.Equal(t, "foo", pod.ObjectMeta.Name)
	assert.Len(t, pod.Spec.Volumes, 0)
	assert.Len(t, pod.Spec.InitContainers, 0)
	require.Len(t, pod.Spec.Containers, 1)
	assert.Equal(t, "nginx", pod.Spec.Containers[0].Name)
	assert.Equal(t, "nginx:1.15.3", pod.Spec.Containers[0].Image)
	assert.Equal(t, k8sv1.RestartPolicyAlways, pod.Spec.RestartPolicy)
}

func TestLoadPodWithDefaults_1c(t *testing.T) {
	t.Skip()
	pod, err := LoadPodWithDefaults([]byte(pod1c_descriptor))
	require.NoError(t, err, "Unable to load pod")
	assert.NotNil(t, pod, "Unable to load pod")
	assert.Equal(t, "foo", pod.ObjectMeta.Name)
	assert.Len(t, pod.Spec.Volumes, 0)
	assert.Len(t, pod.Spec.InitContainers, 0)
	require.Len(t, pod.Spec.Containers, 1)
	assert.Equal(t, "nginx", pod.Spec.Containers[0].Name)
	assert.Equal(t, "nginx:1.15.3", pod.Spec.Containers[0].Image)
	assert.Equal(t, k8sv1.RestartPolicyAlways, pod.Spec.RestartPolicy)
}

func TestLoadPodWithDefaults_1d(t *testing.T) {
	t.Skip()
	pod, err := LoadPodWithDefaults([]byte(pod1d_descriptor))
	require.NoError(t, err, "Unable to load pod")
	assert.NotNil(t, pod, "Unable to load pod")
	assert.Equal(t, "foo", pod.ObjectMeta.Name)
	assert.Len(t, pod.Spec.Volumes, 0)
	assert.Len(t, pod.Spec.InitContainers, 0)
	require.Len(t, pod.Spec.Containers, 2)
	assert.Equal(t, "baz", pod.Spec.Containers[0].Name)
	assert.Equal(t, "nginx:1.15.3", pod.Spec.Containers[0].Image)
	assert.Equal(t, "nginx", pod.Spec.Containers[1].Name)
	assert.Equal(t, "nginx:1.15.3", pod.Spec.Containers[1].Image)
	assert.Equal(t, k8sv1.RestartPolicyNever, pod.Spec.RestartPolicy)
}
