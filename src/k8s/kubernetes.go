package k8s

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os/user"
	"path/filepath"
	"strings"

	v1 "k8s.io/api/core/v1"
	. "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const defaultEmail = "awsregrenew@demo.test"

func GetClient() (*kubernetes.Clientset, error) {
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
	_, err := client.CoreV1().Secrets(namespace).Get(context.Background(), name, GetOptions{})
	if nil != err {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
		return err
	}
	var do = &DeleteOptions{}
	return client.CoreV1().Secrets(namespace).Delete(context.Background(), name, *do)
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

func UpdatePassword(namespace, name, username, password, server string) error {
	client, err := GetClient()
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
	var co = &CreateOptions{}
	_, err = client.CoreV1().Secrets(namespace).Create(context.Background(), secret, *co)
	return err
}
