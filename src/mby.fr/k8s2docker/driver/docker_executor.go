package driver

import (
	"fmt"
	"log"
	"strings"

	"mby.fr/utils/cmdz"
	//"mby.fr/utils/promise"

	"mby.fr/k8s2docker/descriptor"
	"mby.fr/utils/collections"

	corev1 "k8s.io/api/core/v1"
)

type Executor struct {
	translator Translator
	config     []string
	forkCount  int
}

func (e Executor) Apply(namespace string, resource any) (err error) {
	switch res := resource.(type) {
	case corev1.Pod:
		err = e.createPod(namespace, res)
	default:
		err = fmt.Errorf("Cannot Apply type: %T not supported !", res)
	}
	return

}

func (e Executor) ListNamespaces() (namespaces []corev1.Namespace, err error) {
	return e.translator.listNamespaces()
}

func (e Executor) ListPods(namespace string) (pods []corev1.Pod, err error) {
	//return e.translator.listPods(namespace)
	containersMap, err := e.translator.listPodContainers(namespace, "")
	if err != nil {
		return nil, err
	}
	containersNames := collections.Keys(&containersMap)
	podNames := collections.Map(containersNames, func(in string) string {
		// return podName part
		return strings.Split(in, ContainerName_Separator)[1]
	})
	podNames = collections.Deduplicate[string](&podNames)
	for _, podName := range podNames {
		p, err := e.describePod(namespace, podName)
		if err != nil {
			return nil, err
		}
		pods = append(pods, p)
	}
	return
}

func (e Executor) Delete(namespace, kind, name string) (err error) {
	switch kind {
	case "Pod":
		err = e.deletePod(namespace, name)
	default:
		err = fmt.Errorf("Cannot List kind: %s not supported yet !", kind)
	}
	return
}

func (e Executor) describePod(namespace, name string) (corev1.Pod, error) {
	pod := descriptor.BuildDefaultPod(namespace, name)
	containersMap, err := e.translator.listPodContainers(namespace, name)
	if err != nil {
		return corev1.Pod{}, err
	}
	containers := collections.Values(&containersMap)
	pod.Spec.Containers = containers
	descriptor.CompleteK8sResourceDefaults(&pod)
	return pod, nil
}

func (e Executor) deletePod(namespace, name string) (err error) {
	log.Printf("Deleting pod %s in ns %s ...", name, namespace)

	execs, err := e.translator.deletePod(namespace, name)
	if err != nil {
		return err
	}
	_, err = execs.BlockRun()
	if err != nil {
		return err
	}
	//return fmt.Errorf("deletePod not implemented yet !")
	return
}

func (e Executor) createPod(namespace string, pod corev1.Pod) (err error) {
	log.Printf("Creating pod %s in ns %s ...", pod.ObjectMeta.Name, namespace)

	podName := forgeResName(namespace, pod)

	netId, err := e.translator.podNetworkId(namespace, pod)
	if err != nil {
		return err
	}
	if netId == "" {
		exec, err := e.translator.createPodNetwork(namespace, pod)
		if err != nil {
			return err
		}
		_, err = exec.BlockRun()
		if err != nil {
			return fmt.Errorf("Unable to create Network for pod: %s. Caused by: %w", podName, err)
		}
	}

	execs, err := e.translator.createPodRootContainer(namespace, pod)
	if err != nil {
		return err
	}
	_, err = execs.BlockRun()
	if err != nil {
		return fmt.Errorf("Unable to create Root Container for pod: %s. Caused by: %w", podName, err)
	}

	volExec := cmdz.Parallel()
	for _, volume := range pod.Spec.Volumes {
		volId, err := e.translator.volumeId(namespace, volume)
		if err != nil {
			return err
		}
		if volId == "" {
			exec, err := e.translator.createVolume(namespace, volume)
			if err != nil {
				return err
			}
			volExec.Add(exec)
		}
	}
	_, err = volExec.ErrorOnFailure(true).BlockRun()
	if err != nil {
		return fmt.Errorf("Unable to create a volume ! Caused by: %w", err)
	}

	ictExec := cmdz.Parallel()
	for _, container := range pod.Spec.InitContainers {
		ctId, err := e.translator.podContainerId(namespace, pod, container)
		if err != nil {
			return err
		}
		if ctId != "" {
			ctName := podContainerName(namespace, pod, container)
			err = fmt.Errorf("Init Container %s already exists !", ctName)
			return err
		} else {
			exec, err := e.translator.createPodContainer(namespace, pod, container, true)
			if err != nil {
				return err
			}
			ictExec.Add(exec)
		}
	}
	/*
		if len(ictExecs) > 0 {
			err = cmdz.BlockParallel(e.forkCount, ictExecs...)
			if err != nil {
				return fmt.Errorf("Unable to run Init Containers for pod %s. Caused by: %w", podName, err)
			}
		}*/
	_, err = ictExec.ErrorOnFailure(true).BlockRun()
	if err != nil {
		return fmt.Errorf("Unable to create init containers ! Caused by: %w", err)
	}

	ctExec := cmdz.Parallel()
	for _, container := range pod.Spec.Containers {
		ctId, err := e.translator.podContainerId(namespace, pod, container)
		if err != nil {
			return err
		}
		if ctId == "" {
			exec, err := e.translator.createPodContainer(namespace, pod, container, false)
			if err != nil {
				return err
			}
			ctExec.Add(exec)
		}
	}
	/*
		if len(ctExecs) > 0 {
			err = cmdz.BlockParallel(e.forkCount, ctExecs...)
			if err != nil {
				return fmt.Errorf("Unable to run Containers for pod %s. Caused by: %w", podName, err)
			}
		}*/
	_, err = ctExec.ErrorOnFailure(true).BlockRun()
	if err != nil {
		return fmt.Errorf("Unable to create containers ! Caused by: %w", err)
	}

	/*
		ctExec2, err := e.translator.commitPod(namespace, pod)
		if err != nil {
			return fmt.Errorf("Unable to commit pod config ! Caused by: %w", err)
		}
		_, err = ctExec2.FailOnError().BlockRun()
		if err != nil {
			return fmt.Errorf("Unable to commit pod config ! Caused by: %w", err)
		}
	*/
	return
}

func (e Executor) updatePod(namespace string, pod corev1.Pod) (err error) {
	log.Printf("Updating pod %s in ns %s ...", pod.ObjectMeta.Name, namespace)

	// TODO: how to use a kubelet ?
	// Kubelet should be responsible for create / update / delete

	// TODO : check for pod existance. If it already exists and phase is ok check for "updateness"
	// THEN update pod or delete/create pod

	ctId, err := e.translator.podRootContainerId(namespace, pod)
	if err != nil {
		return err
	}
	if ctId == "" {
		// Create pod
	} else {
		podPhase, err := e.translator.getCommitedPodPhase(namespace, pod)
		if err != nil {
			log.Printf("Swallowed error: %s", err)
			// TODO recreate pod
		}
		// TODO what to do with podPhase ?
		_ = podPhase

		/*
			commitedPod, err := e.translator.getCommitedPod(namespace, pod)
			if err != nil {
				log.Printf("Swallowed error: %s", err)
				// TODO recreate pod
			}
			if podPhase != nil && commitedPod != nil {
				diff := compare.ComparePods(*commitedPod, pod)
				if diff.IsUpdatable() {
					// TODO update pod
				}
			}
		*/
	}

	return
}

func DockerExecutor() K8sDriver {
	translator := Translator{"docker"}
	executor := Executor{translator: translator}
	return executor
}
