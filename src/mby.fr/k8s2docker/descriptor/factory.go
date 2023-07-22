package descriptor

import (
	k8sv1 "k8s.io/api/core/v1"
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

func BuildNamespace(name string) (res k8sv1.Namespace) {
	res.TypeMeta = BuildTypeMeta("Namespace", "")
	res.ObjectMeta = BuildObjectMeta("", name)
	return
}

func BuildDefaultPod(namespace, name string) (res k8sv1.Pod) {
	res.TypeMeta = BuildTypeMeta("Pod", "")
	res.ObjectMeta = BuildObjectMeta(namespace, name)
	CompletePodSpecDefaults(&res.Spec)
	return
}

func BuildDefaultContainer(name, image string) (res k8sv1.Container) {
	res.Name = name
	res.Image = image
	CompleteContainerDefaults(&res)
	return
}

func BuildVolumeMount(name, mountPath string) (res k8sv1.VolumeMount) {
	res.Name = name
	res.MountPath = mountPath
	return
}
