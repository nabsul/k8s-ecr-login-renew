package main

import (
	"errors"
	"fmt"
	"github.com/nabsul/k8s-ecr-login-renew/src/aws"
	"github.com/nabsul/k8s-ecr-login-renew/src/k8s"
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
	credentials, err := aws.GetUserAndPass()
	checkErr(err)
	fmt.Printf("Successfully fetched %d docker credentials\n", len(credentials))

	namespaces, err := k8s.GetNamespaces(namespaceList)
	checkErr(err)
	fmt.Printf("Updating kubernetes secret [%s] in %d namespaces\n", name, len(namespaces))

	failed := false
	for _, ns := range namespaces {
		fmt.Printf("Updating secret in namespace [%s]... ", ns)
		err = k8s.UpdatePassword(ns, name, credentials)
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
