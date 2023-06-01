package pod

import (
	k8s "k8s.io/api"
)

type Executor {
	config []string
	translator Translator
}

func (e Executor) exec(commands []string, retries int) (status int, err error) {

}

func (e Executor) Create(resource any) (err error) {
	switch res := resource.(type) {
	case k8s.core.v1.Pod:
		return e.createPod(res)
	}
}

func (e Executor) createPod(pod k8s.core.v1.Pod) (err error) {
	cmds, err := translator.createNetworkOwnerContainer(pod)
	if err != nil {
		return err
	}
	status, err := e.exec(cmds, -1)

	for _, volume := pod.Spec.Volumes {
		cmds, err := translator.createVolume(volume)
		if err != nil {
			return err
		}
		status, err := e.exec(cmds, -1)
	}

	for _, container := pod.Spec.InitContainers {
		cmds, err := translator.createContainer(container)
		if err != nil {
			return err
		}
		status, err := e.exec(cmds, -1)
	}

	for _, container := pod.Spec.Containers {
		cmds, err := translator.createContainer(container)
		if err != nil {
			return err
		}
		status, err := e.exec(cmds, -1)
	}
}

type Translator {
	binary string
}

func (t Translator) createNetworkOwnerContainer(pod k8s.core.v1.Pod) (cmds []string, err error) {
	podName := pod.ObjectMeta.Name
	podMainCtCpus := "0.05"
	podMainCtMemory := "64m"
	networkName := ""
	addHostRules := ""
	pauseImage := ""
	cmd := fmt.Sprintf(	"docker run -d --name '%s' --cpus=%s --memory=%s --memory-swap=%s --memory-swappiness=0" + 
						" --network '%s' %s --restart=always '%s' /bin/sleep inf",
						podName, podMainCtCpus, podMainCtMemory, podMainCtMemory, networkName, addHostRules, pauseImage
	)
	append(cmds, cmd)
}

func (t Translator) createVolume(vol k8s.core.v1.Volume) (cmds []string, err error) {

}

func (t Translator) createContainer(vol k8s.core.v1.Container) (cmds []string, err error) {
	
}
