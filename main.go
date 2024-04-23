package main

import (
	"log"
	"os"
	"time"

	"github.com/sikalabs/github-apps-pull-secret-sync/pkg/ghcr"
	"github.com/sikalabs/github-apps-pull-secret-sync/pkg/kubernetes"
	"github.com/sikalabs/github-apps-pull-secret-sync/version"
)

func main() {
	if len(os.Args) < 5 {
		log.Fatalf("Usage github-apps-pull-secret-sync "+version.Version+": %s <githubAppID> <githubInstallationID> <privateKeyPath> <username> [<namespace> ...]", os.Args[0])
	}

	log.Println("Starting github-apps-pull-secret-sync", version.Version)

	for {
		token := ghcr.GetGhcrToken(os.Args[1], os.Args[2], os.Args[3])
		dockerConfigJson := ghcr.CreateDockerConfigJson(os.Args[4], token)
		namespaces := os.Args[5:]
		if len(namespaces) == 0 {
			namespaces = kubernetes.GetNamespaces()
		}
		for _, namespace := range namespaces {
			kubernetes.CreateOrUpdareSecretDockerConfigJson("github-apps-pull-secret", namespace, dockerConfigJson)
		}
		log.Println("Sleeping 50 minutes ...")
		time.Sleep(50 * time.Minute)
	}
}
