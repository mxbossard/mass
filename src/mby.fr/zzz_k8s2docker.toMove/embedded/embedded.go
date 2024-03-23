package main

import (
	"log"
	"time"

	"mby.fr/k8s2docker/daemon"
	"mby.fr/k8s2docker/repo"
	"mby.fr/k8s2docker/server"
)

func main() {
	dbPath := "/tmp/k8s2dockerDb"
	tcpAddr := ":8080"

	log.Printf("Starting k8s2docker daemon and HTTP server ...")

	repo.InitDb(dbPath)

	// Not blocking
	daemon.Start(2*time.Second, 20*time.Second)
	defer daemon.Stop()

	// Blocking
	server.Start(tcpAddr)

}
