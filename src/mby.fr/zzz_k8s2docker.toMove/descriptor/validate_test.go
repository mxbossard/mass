package descriptor

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	typeMeta_empty = metav1.TypeMeta{}
	typeMeta_part  = metav1.TypeMeta{Kind: "Pod"}
	typeMeta_valid = metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"}

	objectMeta_empty = metav1.ObjectMeta{}
	objectMeta_part  = metav1.ObjectMeta{Namespace: "foo"}
	objectMeta_valid = metav1.ObjectMeta{Name: "bar"}

	namespace_empty  = corev1.Namespace{}
	namespace_part_1 = corev1.Namespace{TypeMeta: typeMeta_valid}
	namespace_part_2 = corev1.Namespace{ObjectMeta: objectMeta_valid}
	namespace_valid  = corev1.Namespace{TypeMeta: typeMeta_valid, ObjectMeta: objectMeta_valid}

	container_empty  = corev1.Container{}
	container_part_1 = corev1.Container{Name: "foo"}
	container_part_2 = corev1.Container{Image: "alpine"}
	container_part_3 = corev1.Container{Name: "foo", Image: "alpine"}
	container_part_4 = corev1.Container{Name: "foo", ImagePullPolicy: corev1.PullNever}
	container_valid  = corev1.Container{Name: "foo", Image: "alpine", ImagePullPolicy: corev1.PullNever}

	volume_empty = corev1.Volume{}
	volume_valid = corev1.Volume{Name: "foo"}

	volumeMount_empty = corev1.VolumeMount{}
	volumeMount_part  = corev1.VolumeMount{Name: "foo"}
	volumeMount_valid = corev1.VolumeMount{Name: "foo", MountPath: "bar"}

	podspec_empty = corev1.PodSpec{}
	podspec_part  = corev1.PodSpec{Containers: []corev1.Container{container_valid}}
	podspec_valid = corev1.PodSpec{RestartPolicy: corev1.RestartPolicyNever}

	pod_empty  = corev1.Pod{}
	pod_part_1 = corev1.Pod{TypeMeta: typeMeta_valid}
	pod_part_2 = corev1.Pod{ObjectMeta: objectMeta_valid}
	pod_part_3 = corev1.Pod{Spec: podspec_valid}
	pod_part_4 = corev1.Pod{TypeMeta: typeMeta_valid, ObjectMeta: objectMeta_valid}
	pod_valid  = corev1.Pod{TypeMeta: typeMeta_valid, ObjectMeta: objectMeta_valid, Spec: podspec_valid}
)

func TestValidateTypeMeta(t *testing.T) {
	aggErrors := ValidateTypeMeta(typeMeta_empty, "Pod", "p0")
	require.NotNil(t, aggErrors)
	assert.Error(t, aggErrors.Return())

	aggErrors = ValidateTypeMeta(typeMeta_part, "Pod", "p1")
	require.NotNil(t, aggErrors)
	assert.Error(t, aggErrors.Return())

	aggErrors = ValidateTypeMeta(typeMeta_valid, "Pod", "pValid")
	require.NotNil(t, aggErrors)
	assert.NoError(t, aggErrors.Return())
}

func TestValidateObjectMeta(t *testing.T) {
	aggErrors := ValidateObjectMeta(objectMeta_empty, "Pod", "p0")
	require.NotNil(t, aggErrors)
	assert.Error(t, aggErrors.Return())

	aggErrors = ValidateObjectMeta(objectMeta_part, "Pod", "p1")
	require.NotNil(t, aggErrors)
	assert.Error(t, aggErrors.Return())

	aggErrors = ValidateObjectMeta(objectMeta_valid, "Pod", "pValid")
	require.NotNil(t, aggErrors)
	assert.NoError(t, aggErrors.Return())
}

func TestValidateNamespace(t *testing.T) {
	err := ValidateNamespace(namespace_empty, "ns0")
	assert.Error(t, err)

	err = ValidateNamespace(namespace_part_1, "ns1")
	assert.Error(t, err)

	err = ValidateNamespace(namespace_part_2, "ns2")
	assert.Error(t, err)

	err = ValidateNamespace(namespace_valid, "nsValid")
	assert.NoError(t, err)
}

func TestValidateK8sResource_Namespace(t *testing.T) {
	_, err := ValidateK8sResource(namespace_empty, "ns0")
	assert.Error(t, err)

	_, err = ValidateK8sResource(namespace_part_1, "ns1")
	assert.Error(t, err)

	_, err = ValidateK8sResource(namespace_part_2, "ns2")
	assert.Error(t, err)

	_, err = ValidateK8sResource(namespace_valid, "nsValid")
	assert.NoError(t, err)
}

func TestValidateMappedK8sResource_Namespace(t *testing.T) {
	emptyNs := map[string]any{}
	partialNs := map[string]any{"kind": "Namespace", "apiVersion": "v1"}
	validNs := map[string]any{"kind": "Namespace", "apiVersion": "v1", "metadata": map[string]any{"name": "foo"}}

	var err error
	_, err = ValidateMappedK8sResource(emptyNs, "Namespace", "ns1")
	assert.Error(t, err)

	_, err = ValidateMappedK8sResource(partialNs, "Namespace", "ns2")
	assert.Error(t, err)

	_, err = ValidateMappedK8sResource(validNs, "Namespace", "ns3")
	assert.NoError(t, err)
}

func TestValidateVolumeMount(t *testing.T) {
	aggErrors := ValidateVolumeMount(volumeMount_empty, "spec.containers[0].volumeMounts[0]", "Pod", "p0")
	require.NotNil(t, aggErrors)
	assert.Error(t, aggErrors.Return())

	aggErrors = ValidateVolumeMount(volumeMount_part, "spec.containers[0].volumeMounts[0]", "Pod", "p1")
	require.NotNil(t, aggErrors)
	assert.Error(t, aggErrors.Return())

	aggErrors = ValidateVolumeMount(volumeMount_valid, "spec.containers[0].volumeMounts[0]", "Pod", "pValid")
	require.NotNil(t, aggErrors)
	assert.NoError(t, aggErrors.Return())
}

func TestValidateContainer(t *testing.T) {
	aggErrors := ValidateContainer(container_empty, "spec.containers[0]", "Pod", "p0")
	require.NotNil(t, aggErrors)
	assert.Error(t, aggErrors.Return())

	aggErrors = ValidateContainer(container_part_1, "spec.containers[0]", "Pod", "p1")
	require.NotNil(t, aggErrors)
	assert.Error(t, aggErrors.Return())

	aggErrors = ValidateContainer(container_part_2, "spec.containers[0]", "Pod", "p2")
	require.NotNil(t, aggErrors)
	assert.Error(t, aggErrors.Return())

	aggErrors = ValidateContainer(container_part_3, "spec.containers[0]", "Pod", "p3")
	require.NotNil(t, aggErrors)
	assert.Error(t, aggErrors.Return())

	aggErrors = ValidateContainer(container_part_4, "spec.containers[0]", "Pod", "p4")
	require.NotNil(t, aggErrors)
	assert.Error(t, aggErrors.Return())

	aggErrors = ValidateContainer(container_valid, "spec.containers[0]", "Pod", "pValid")
	require.NotNil(t, aggErrors)
	assert.NoError(t, aggErrors.Return())
}

func TestValidatePodSpec(t *testing.T) {
	aggErrors := ValidatePodSpec(podspec_empty, "spec", "Pod", "p0")
	require.NotNil(t, aggErrors)
	assert.Error(t, aggErrors.Return())

	aggErrors = ValidatePodSpec(podspec_part, "spec", "Pod", "p1")
	require.NotNil(t, aggErrors)
	assert.Error(t, aggErrors.Return())

	aggErrors = ValidatePodSpec(podspec_valid, "spec", "Pod", "pValid")
	require.NotNil(t, aggErrors)
	assert.NoError(t, aggErrors.Return())

}

func TestValidatePod(t *testing.T) {
	err := ValidatePod(pod_empty, "p0")
	assert.Error(t, err)

	err = ValidatePod(pod_part_1, "p1")
	assert.Error(t, err)

	err = ValidatePod(pod_part_2, "p2")
	assert.Error(t, err)

	err = ValidatePod(pod_part_3, "p3")
	assert.Error(t, err)

	err = ValidatePod(pod_part_4, "p4")
	assert.Error(t, err)

	err = ValidatePod(pod_valid, "pValid")
	assert.NoError(t, err)

}

func TestValidateK8sResource_Pod(t *testing.T) {
	_, err := ValidateK8sResource(pod_empty, "p0")
	assert.Error(t, err)

	_, err = ValidateK8sResource(pod_part_1, "p1")
	assert.Error(t, err)

	_, err = ValidateK8sResource(pod_part_2, "p2")
	assert.Error(t, err)

	_, err = ValidateK8sResource(pod_part_3, "p3")
	assert.Error(t, err)

	// Should not fail because the resource is completed with defaults and pass validation
	_, err = ValidateK8sResource(pod_part_4, "p4")
	assert.NoError(t, err)

	_, err = ValidateK8sResource(pod_valid, "pValid")
	assert.NoError(t, err)
}

func TestValidateMappedK8sResource_Pod(t *testing.T) {
	emptyPod_invalid := map[string]any{}
	partialPod1_invalid := map[string]any{"kind": "Pod", "apiVersion": "v1"}
	partialPod2_valid := map[string]any{"kind": "Pod", "apiVersion": "v1",
		"metadata": map[string]any{"name": "foo"},
	}
	partialPod3_invalid := map[string]any{"kind": "Pod", "apiVersion": "v1",
		"metadata": map[string]any{"name": "foo"},
		"spec": map[string]any{
			"containers": map[string]any{
				"name": "foo",
			},
		},
	}
	partialPod4_valid := map[string]any{"kind": "Pod", "apiVersion": "v1",
		"metadata": map[string]any{"name": "foo"},
		"spec": map[string]any{
			"containers": []map[string]any{
				map[string]any{
					"name":  "foo",
					"image": "alpine",
				},
			},
		},
	}
	fullPod1_valid := map[string]any{"kind": "Pod", "apiVersion": "v1",
		"metadata": map[string]any{"name": "foo"},
		"spec": map[string]any{
			"containers": []map[string]any{
				map[string]any{
					"name":            "foo",
					"image":           "alpine",
					"imagePullPolicy": "Never",
				},
			},
		},
	}

	var err error
	_, err = ValidateMappedK8sResource(emptyPod_invalid, "Pod", "p0")
	assert.Error(t, err)

	_, err = ValidateMappedK8sResource(partialPod1_invalid, "Pod", "p1")
	assert.Error(t, err)

	_, err = ValidateMappedK8sResource(partialPod2_valid, "Pod", "p2")
	assert.NoError(t, err)

	_, err = ValidateMappedK8sResource(partialPod3_invalid, "Pod", "p3")
	assert.Error(t, err)

	_, err = ValidateMappedK8sResource(partialPod4_valid, "Pod", "p4")
	require.NoError(t, err)

	_, err = ValidateMappedK8sResource(fullPod1_valid, "Pod", "pFull")
	assert.NoError(t, err)
}
