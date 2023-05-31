package resources

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	//k8s "k8s.io/api"
	"os"
	"fmt"
	"testing"
	"path/filepath"

	//"mby.fr/mass/internal/commontest"
	"mby.fr/utils/test"
)


func TestReadService(t *testing.T) {
	expectedProjectDir, err := test.BuildRandTempPath()
	os.MkdirAll(expectedProjectDir, 0755)
	defer os.RemoveAll(expectedProjectDir)

	podSpecFilePath := filepath.Join(expectedProjectDir, "my-pod.yaml")
	expectedServiceName := "my-service"
	expectedServiceResFileName := forgeServiceResFilename(expectedServiceName)
	expectedServiceResFilePath := filepath.Join(expectedProjectDir, expectedServiceResFileName)
	expectedServiceResContent := fmt.Sprintf(`
resourceKind: service
ServiceType: K8sService
ServiceFile: %s
`, podSpecFilePath)
	err = os.WriteFile(expectedServiceResFilePath, []byte(expectedServiceResContent), 0644)
	require.NoError(t, err, "should not error")

	podSpec := `
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
	err = os.WriteFile(podSpecFilePath, []byte(podSpec), 0644)
	require.NoError(t, err, "should not error")

	service, err := Read[Pod](expectedServiceResFilePath)
	require.NoError(t, err, "should not error")
	
	assert.Equal(t, expectedServiceName, service.Name())
	assert.Equal(t, expectedProjectDir, service.Dir())
}