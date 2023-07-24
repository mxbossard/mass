package driver

import (

	//"mby.fr/utils/promise"

	"fmt"

	corev1 "k8s.io/api/core/v1"
)

type DockerDriver struct {
	translater Translater
	config     []string
	forkCount  int
}

func (d DockerDriver) ApplyNamespace(in corev1.Namespace) (out corev1.Namespace, err error) {
	err = fmt.Errorf("driver.ApplyNamespace() not implemented yet !")
	return
}

func (d DockerDriver) DeleteNamespace(in corev1.Namespace) (out corev1.Namespace, err error) {
	err = fmt.Errorf("driver.DeleteNamespace() not implemented yet !")
	return
}

func (d DockerDriver) DescribeNamespace(namespace string) (out corev1.Namespace, err error) {
	err = fmt.Errorf("driver.DescribeNamespace() not implemented yet !")
	return
}

func (d DockerDriver) ListNamespaces() (namespaces []corev1.Namespace, err error) {
	err = fmt.Errorf("driver.ListNamespaces() not implemented yet !")
	return
}

func (d DockerDriver) ApplyPod(in corev1.Pod) (out corev1.Pod, err error) {
	err = fmt.Errorf("driver.ApplyPod() not implemented yet !")
	return
}

func (d DockerDriver) DeletePod(in corev1.Pod) (out corev1.Pod, err error) {
	err = fmt.Errorf("driver.DeletePod() not implemented yet !")
	return
}

func (d DockerDriver) DescribePod(namespace, name string) (out corev1.Pod, err error) {
	err = fmt.Errorf("driver.DescribePod() not implemented yet !")
	return
}

func (d DockerDriver) ListPods(namespace string) (pods []corev1.Pod, err error) {
	err = fmt.Errorf("driver.ListPods() not implemented yet !")
	return
}

func NewDockerDriver() DockerDriver {
	translater := DockerTranslater{"docker"}
	driver := DockerDriver{translater: translater}
	return driver
}
