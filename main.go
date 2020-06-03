package main

import (
	"fmt"
	"os"
	"time"
)

const awsSecretName = "DOCKER_SECRET_NAME"
const targetNamespaceName = "TARGET_NAMESPACE"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Running at " + time.Now().UTC().String())

	name := os.Getenv(awsSecretName)
	if name == "" {
		panic(fmt.Sprintf("Environment variable %s is required", name))
	}

	targetNamespace := os.Getenv(targetNamespaceName)
	if targetNamespace == "" {
		targetNamespace = "default"
	}

	fmt.Print("Fetching auth data from AWS... ")
	username, password, server := getUserAndPass()
	fmt.Println("Success.")

	fmt.Printf("Updating kubernetes secret [%s]... ", name)

	for _, ns := range getNamespaces(targetNamespace) {
		updatePassword(name, username, password, server, ns)
	}

	fmt.Println("Success.")
}
