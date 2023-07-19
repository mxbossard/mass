package driver

import (
	corev1 "k8s.io/api/core/v1"
)

type K8sDriver interface {
	Apply(string, any) (err error)
	Delete(string, string, string) (err error)
	ListPods(string) ([]corev1.Pod, error)
}
