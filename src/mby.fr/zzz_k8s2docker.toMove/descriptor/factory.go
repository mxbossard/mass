package descriptor

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func BuildTypeMeta(kind, apiVersion string) (typemeta metav1.TypeMeta) {
	typemeta.Kind = kind
	if apiVersion == "" {
		apiVersion = "v1"
	}
	typemeta.APIVersion = apiVersion
	return
}

func BuildObjectMeta(namespace, name string) (metadata metav1.ObjectMeta) {
	metadata.Name = name
	metadata.Namespace = namespace
	return
}

func BuildNamespace(name string) (res corev1.Namespace) {
	res.TypeMeta = BuildTypeMeta("Namespace", "")
	res.ObjectMeta = BuildObjectMeta("", name)
	return
}

func BuildDefaultPod(namespace, name string) (res corev1.Pod) {
	res.TypeMeta = BuildTypeMeta("Pod", "")
	res.ObjectMeta = BuildObjectMeta(namespace, name)
	CompletePodSpecDefaults(&res.Spec)
	return
}

func BuildDefaultContainer(name, image string) (res corev1.Container) {
	res.Name = name
	res.Image = image
	CompleteContainerDefaults(&res)
	return
}

func BuildVolumeMount(name, mountPath string) (res corev1.VolumeMount) {
	res.Name = name
	res.MountPath = mountPath
	return
}

func BuildDefaultHostPathVolume(name, hostPath string) (res corev1.Volume) {
	res.Name = name
	defaultType := corev1.HostPathUnset
	res.VolumeSource.HostPath = &corev1.HostPathVolumeSource{
		Path: hostPath,
		Type: &defaultType,
	}
	return
}

func BuildDefaultEmptyDirVolume(name string) (res corev1.Volume) {
	res.Name = name
	res.VolumeSource.EmptyDir = &corev1.EmptyDirVolumeSource{
		Medium:    corev1.StorageMediumDefault,
		SizeLimit: nil,
	}
	return
}

func BuildDefaultSecurityContext(runAsUser, runAsGroup *int64) (res corev1.SecurityContext) {
	res.Privileged = boolPtr(false)
	res.ReadOnlyRootFilesystem = boolPtr(false)
	res.RunAsNonRoot = boolPtr(false)
	res.RunAsUser = runAsUser
	res.RunAsGroup = runAsGroup
	return
}

func boolPtr(in bool) *bool {
	return &in
}
