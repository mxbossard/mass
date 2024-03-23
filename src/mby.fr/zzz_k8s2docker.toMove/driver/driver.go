package driver

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"mby.fr/utils/cmdz"
)

type K8sDriver interface {
	Apply(string, any) (err error)
	Delete(string, string, string) (err error)
	ListNamespaces() ([]corev1.Namespace, error)
	ListPods(string) ([]corev1.Pod, error)
}

type K8sDriver3 interface {
	ApplyNamespace(corev1.Namespace) (corev1.Namespace, error)
	DeleteNamespace(corev1.Namespace) (corev1.Namespace, error)
	DescribeNamespace(namespace string) (corev1.Namespace, error)
	ListNamespaces() ([]corev1.Namespace, error)

	ApplyPod(corev1.Pod) (corev1.Pod, error)
	DeletePod(corev1.Pod) (corev1.Pod, error)
	DescribePod(namespace, name string) (corev1.Pod, error)
	ListPods(namespace string) ([]corev1.Pod, error)
}

// Driver Helper funcs

// Return a description of existing resource of name & kind in namespace
// Return nil if not exists
func Describe[T any](driver K8sDriver3, namespace, kind, name string) (resource T, err error) {
	err = fmt.Errorf("driver.Describe() does not exists yet")
	return
}

// Return a list of existing resources of kind in namespace
func List[T any](namespace string) (resources []T, err error) {
	err = fmt.Errorf("driver.List() does not exists yet")
	return
}

// Delete existing resource of name & kind in namespace
func Delete[T any](namespace, name string) (resource T, err error) {
	err = fmt.Errorf("driver.Delete() does not exists yet")
	return
}

// Compare existing resource in namespace with supplied one
// Create it if does not exist
// Update it if does exist
func Apply[T any](namespace, kind string, in T) (out T, err error) {
	err = fmt.Errorf("driver.Apply() does not exists yet")
	return
}

// Compare all existing resources of kind in namespace with suplied ones
// Create all not existing resources
// Update all existing resources
// Delete all not listed resources
func ApplyAll[T any](namespace, kind string, in ...T) (out []T, err error) {
	err = fmt.Errorf("driver.ApplyAll() does not exists yet")
	return
}

type Translater interface {
	CreateNamespace(corev1.Namespace) cmdz.Executer
	UpdateNamespace(corev1.Namespace) cmdz.Executer
	DeleteNamsepace(namespace string) cmdz.Executer
	ListNamespaceNames() cmdz.Formatter[[]string]
	//DescribeNamespace(string) (corev1.Namespace, error)
	//ListNamespaces() ([]corev1.Namespace, error)

	//VolumeId(string, corev1.Volume) (cmdz.Executer, error)
	CreateVolume(namespace, podName string, vol corev1.Volume) (cmdz.Executer, error)
	DeleteVolume(namepsace, podName, name string) cmdz.Executer
	DescribeVolume(namepsace, podName, name string) cmdz.Formatter[[]corev1.Volume]
	ListVolumeNames(namespace, podName string) cmdz.Formatter[[]string]
	//DescribeVolume(string, string) (corev1.Volume, error)
	//ListVolumes(string) ([]corev1.Volume, error)

	// Merge Network into Root Container
	//CreatePodNetwork(string, string) (cmdz.Executer, error)
	//PodNetworkId(string, string) (cmdz.Executer, error)
	//PodRootContainerId(string, corev1.Pod) (cmdz.Executer, error)
	//CreatePodRootContainer(namespace string, pod corev1.Pod) cmdz.Executer
	//DeletePodRootContainer(namespace, name string) cmdz.Executer
	//PodContainerId(string, corev1.Pod, corev1.Container) (cmdz.Executer, error)
	SetupPod(namespace string, pod corev1.Pod) cmdz.Executer
	DeletePod(namespace, name string) cmdz.Executer
	DescribePod(namepsace, podName string) cmdz.Formatter[*corev1.Pod]

	CreatePodContainer(namespace string, pod corev1.Pod, ct corev1.Container, init bool) (cmdz.Executer, error)
	UpdatePodContainer(namespace string, pod corev1.Pod, ct corev1.Container) (cmdz.Executer, error)
	DeletePodContainer(namepsace, podName, name string) cmdz.Executer
	DescribePodContainer(namepsace, podName, name string) cmdz.Formatter[[]corev1.Container]
	ListPodContainerNames(namepsace, podName string) cmdz.Formatter[[]string]
	//DescribePodContainer(string, string, string) (corev1.Container, error)
	//ListPodContainers(string, string) ([]corev1.Container, error)

	//ListPodContainers(string, string) (map[string]corev1.Container, error)
	//DeletePodContainer(string, string) (cmdz.Executer, error)
}