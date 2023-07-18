package daemon

import (
	"log"
	"time"

	"mby.fr/k8s2docker/repo"
)

func Start() {
	for {
		err := process()
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
		time.Sleep(1 * time.Second)
	}
}

func process() (err error) {
	namespaces, err := repo.ListNamespaces()
	if err != nil {
		return err
	}

	for _, ns := range namespaces {

	}
}
