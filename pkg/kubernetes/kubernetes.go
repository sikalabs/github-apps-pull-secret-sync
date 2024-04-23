package kubernetes

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateOrUpdareSecretDockerConfigJson(name, namespace, dockerConfigJson string) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
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
			Type: corev1.SecretTypeDockerConfigJson,
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
