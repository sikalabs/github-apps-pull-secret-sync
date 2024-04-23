package kubernetes

import (
	"context"
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
		log.Printf("Secret %s/%s updated\n", namespace, name)
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
		log.Printf("Secret %s/%s created\n", namespace, name)
	}
}

func GetNamespaces() []string {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %s", err.Error())
	}
	namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error getting namespaces: %s", err.Error())
	}
	namespaces := []string{}
	for _, namespace := range namespaceList.Items {
		namespaces = append(namespaces, namespace.Name)
	}
	return namespaces
}
