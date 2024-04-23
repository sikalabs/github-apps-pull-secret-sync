package kubernetes

import (
	"context"
	"fmt"
	"log"

	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func CreateOrUpdareSecretDockerConfigJson(name, namespace, dockerConfigJson string) {
	// Set the kubeconfig to use the in-cluster config.
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfigPath := filepath.Join(home, ".kube", "config")
		kubeconfig = &kubeconfigPath
	}

	// Build the configuration using the kubeconfig.
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	// Create the Kubernetes client.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %s", err.Error())
	}

	// Check if the secret already exists.
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err == nil && secret != nil {
		// Update secret if it exists.
		secret.Type = corev1.SecretTypeDockerConfigJson
		secret.StringData = map[string]string{
			corev1.DockerConfigJsonKey: dockerConfigJson,
		}
		_, updateErr := clientset.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
		if updateErr != nil {
			log.Fatalf("Failed to update secret: %s", updateErr.Error())
		}
		fmt.Println("Secret updated")
	} else {
		// Create secret if it does not exist.
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			StringData: map[string]string{
				corev1.DockerConfigJsonKey: dockerConfigJson,
			},
		}
		_, createErr := clientset.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
		if createErr != nil {
			log.Fatalf("Failed to create secret: %s", createErr.Error())
		}
		fmt.Println("Secret created")
	}
}
