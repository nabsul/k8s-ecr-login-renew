package main

import (
	"encoding/base64"
	"encoding/json"
	"k8s.io/api/core/v1"
	. "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os/user"
	"path/filepath"
)

const defaultEmail = "awsregrenew@demo.test"

func getClient() (*kubernetes.Clientset, error) {
	config, err := getClientConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func getClientConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	u, err := user.Current()
	if nil != err {
		return nil, err
	}

	return clientcmd.BuildConfigFromFlags("", filepath.Join(u.HomeDir, ".kube", "config"))
}

func deleteOldSecret(client *kubernetes.Clientset, name, namespace string) error {
	_, err := client.CoreV1().Secrets(namespace).Get(name, GetOptions{})
	if nil != err {
		return err
	}

	return client.CoreV1().Secrets(namespace).Delete(name, &DeleteOptions{})
}

func createSecret(name, username, password, server string) (*v1.Secret, error) {
	config := map[string]map[string]map[string]string{
		"auths": {
			server: {
				"username": username,
				"password": password,
				"email":    defaultEmail,
				"auth":     base64.StdEncoding.EncodeToString([]byte(username + ":" + password)),
			},
		},
	}

	configJson, err := json.Marshal(config)
	if nil != err {
		return nil, err
	}

	secret := v1.Secret{}
	secret.Name = name
	secret.Type = v1.SecretTypeDockerConfigJson
	secret.Data = map[string][]byte{}
	secret.Data[v1.DockerConfigJsonKey] = configJson
	return &secret, nil
}

func updatePassword(name, username, password, server, namespace string) error {
	client, err := getClient()
	if nil != err {
		return err
	}

	err = deleteOldSecret(client, name, namespace)
	if nil != err {
		return err
	}

	secret, err := createSecret(name, username, password, server)
	if nil != err {
		return err
	}

	_, err = client.CoreV1().Secrets(namespace).Create(secret)
	return err
}
