package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	envVarAwsSecret       = "DOCKER_SECRET_NAME"
	envVarTargetNamespace = "TARGET_NAMESPACE"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Running at " + time.Now().UTC().String())

	name := os.Getenv(envVarAwsSecret)
	if name == "" {
		panic(fmt.Sprintf("Environment variable %s is required", name))
	}

	namespaceList := os.Getenv(envVarTargetNamespace)
	if namespaceList == "" {
		namespaceList = "default"
	}

	fmt.Print("Fetching auth data from AWS... ")
	username, password, server, err := getUserAndPass()
	checkErr(err)
	fmt.Println("Success.")

	fmt.Printf("Updating kubernetes secret [%s]... ", name)

	namespaces, err := getNamespaces(namespaceList)
	checkErr(err)

	failed := false
	for _, ns := range namespaces {
		fmt.Printf("Updating secret [%s] in namespace [%s]... ", name, ns)
		err = updatePassword(name, username, password, server, ns)
		if nil != err {
			fmt.Printf("failed: %s\n", err)
			failed = true
		} else {
			fmt.Println("success")
		}
	}

	if failed {
		panic(errors.New("failed to create one of more Docker login secrets"))
	}

	fmt.Println("Job complete.")
}
