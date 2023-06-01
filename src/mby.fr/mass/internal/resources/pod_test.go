package resources

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	//k8s "k8s.io/api"
	"os"
	"path/filepath"
	"testing"

	//"mby.fr/mass/internal/commontest"
	"mby.fr/utils/test"
)

func TestReadPod(t *testing.T) {
	t.Skip("WIP")
	expectedProjectDir, err := test.BuildRandTempPath()
	os.MkdirAll(expectedProjectDir, 0755)
	defer os.RemoveAll(expectedProjectDir)

	expectedPodName := "my-res"
	expectedPodResFileName := "pod-" + expectedPodName + " .yaml"
	expectedPodResFilePath := filepath.Join(expectedProjectDir, expectedPodResFileName)
	expectedPodResContent := `
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
`
	err = os.WriteFile(expectedPodResFilePath, []byte(expectedPodResContent), 0644)
	require.NoError(t, err, "should not error")

	pod, err := Read[Pod](expectedPodResFilePath)
	require.NoError(t, err, "should not error")

	assert.Equal(t, expectedPodName, pod.Name())
	assert.Equal(t, expectedProjectDir, pod.Dir())
}
