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

func getEnv(name string) string {
	result := os.Getenv(name)
	if result == "" {
		panic("Environment variable " + name + " is required")
	}
	return result
}

func main() {
	fmt.Println("Running at " + time.Now().UTC().String())

	name := getEnv(awsSecretName)
	targetNamespace := getEnv(targetNamespaceName)
	if targetNamespace == ""{
		targetNamespace = "default"
	}

	fmt.Print("Fetching auth data from AWS... ")
	username, password, server := getUserAndPass()
	fmt.Println("Success.")

	fmt.Printf("Updating kubernetes secret [%s]... ", name)
	updatePassword(name, username, password, server, targetNamespace)
	fmt.Println("Success.")
}
