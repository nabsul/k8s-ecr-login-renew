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
	"strings"
)

const defaultEmail = "awsregrenew@demo.test"

func getClient() kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		u, err := user.Current()
		checkErr(err)

		config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(u.HomeDir, ".kube", "config"))
		checkErr(err)
	}

	client, err := kubernetes.NewForConfig(config)
	checkErr(err)

	return *client
}

func deleteOldSecret(client kubernetes.Clientset, name, namespace string) {
	_, err := client.CoreV1().Secrets(namespace).Get(name, GetOptions{})
	if err == nil {
		err = client.CoreV1().Secrets(namespace).Delete(name, &DeleteOptions{})
		checkErr(err)
	}
}

func createSecret(name, username, password, server string) v1.Secret {
	config := map[string]map[string]map[string]string {
		"auths": {
			server: {
				"username": username,
				"password": password,
				"email": defaultEmail,
				"auth": base64.StdEncoding.EncodeToString([]byte(username + ":" + password)),
			},
		},
	}

	configJson, err := json.Marshal(config)
	checkErr(err)

	secret := v1.Secret{}
	secret.Name = name
	secret.Type = v1.SecretTypeDockerConfigJson
	secret.Data = map[string][]byte{}
	secret.Data[v1.DockerConfigJsonKey] = configJson
	return secret
}

func updatePassword(name, username, password, server, namespace string) {
	client := getClient()

	deleteOldSecret(client, name, namespace)

	secret := createSecret(name, username, password, server)

	_, err := client.CoreV1().Secrets(namespace).Create(&secret)
	checkErr(err)
}

func getNamespaces(envVar string) []string {
	list := strings.Split(envVar, ",")
	return list
}

func getAllNamespaces() []string {
	var result []string

	client := getClient()
	opts := ListOptions{}
	first := true

	for first || opts.Continue != "" {
		first = false
		res, err := client.CoreV1().Namespaces().List(opts)
		checkErr(err)
		opts.Continue = res.Continue
		newNames := make([]string, len(res.Items))
		for i, item := range res.Items {
			newNames[i] = item.Name
		}

		result = append(result, newNames...)
	}

	return result
}
