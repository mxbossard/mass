package resources

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8s "k8s.io/api"
	"os"
	"testing"

	"mby.fr/mass/internal/commontest"
)

func TestReadPod(t *testing.T) {
	expectedProjectDir, err := test.BuildRandTempPath()
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	expectedPodName := "my-res"
	expectedPodResFileName := strings.Join("pod-", expectedPodName, " .yaml")
	expectedPodResFilePath := filepath.Join(expectedProjectDir, expectedPodResFileName)
	expectedPodResContent := """
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
"""
	err = os.WriteFile(expectedPodResFilePath, []byte(expectedPodResContent), 0644)
	require.NoError(t, err, "should not error")

	pod, err := Read[Pod](expectedProjectDir, expectedPodName)
	require.NoError(t, err, "should not error")
	

}