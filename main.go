package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/nabsul/k8s-ecr-login-renew/src/aws"
	"github.com/nabsul/k8s-ecr-login-renew/src/k8s"
	"github.com/nabsul/k8s-ecr-login-renew/src/vault"
)

const (
	envVarAwsSecret       = "DOCKER_SECRET_NAME"
	envVarTargetNamespace = "TARGET_NAMESPACE"
	envVarVaultEnable     = "VAULT_ENBALE"
	envVarVaultAddr       = "VAULT_ADDR"
	envVarVaultToken      = "VAULT_TOKEN"
	envVarVaultPath       = "VAULT_SECRET_PATH"
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
		panic(fmt.Sprintf("Environment variable %s is required", envVarAwsSecret))
	}

	namespaceList := os.Getenv(envVarTargetNamespace)
	if namespaceList == "" {
		namespaceList = "default"
	}

	// AWS ECR login variables
	var username, password, server string

	vaultEnable := os.Getenv(envVarVaultEnable)
	if vaultEnable != "true" {

		vaultAddr := os.Getenv(envVarVaultAddr)
		if vaultAddr == "" {
			panic(fmt.Sprintf("Environment variable %s is required", envVarVaultAddr))
		}

		vaultToken := os.Getenv(envVarVaultToken)
		if vaultToken == "" {
			panic(fmt.Sprintf("Environment variable %s is required", envVarVaultToken))
		}

		vaultPath := os.Getenv(envVarVaultPath)
		if vaultPath == "" {
			panic(fmt.Sprintf("Environment variable %s is required", envVarVaultPath))
		}
		result, err := vault.LoginVault(vaultAddr, vaultToken, vaultPath)

		fmt.Print("Fetching auth data from AWS Securely... ")
		username, password, server, err = aws.GetUserAndPass(result["AWS_ACCESS_KEY_ID"], result["AWS_SECRET_ACCESS_KEY"], result["AWS_REGION"])
		checkErr(err)
		fmt.Println("Success.", username, password, server)

	} else {
		fmt.Print("Fetching auth data from AWS... ")
		var err error
		username, password, server, err = aws.GetUserAndPass("", "", "")
		checkErr(err)
		fmt.Println("Success.")
	}

	namespaces, err := k8s.GetNamespaces(namespaceList)
	checkErr(err)
	fmt.Printf("Updating kubernetes secret [%s] in %d namespaces\n", name, len(namespaces))

	failed := false
	for _, ns := range namespaces {
		fmt.Printf("Updating secret in namespace [%s]... ", ns)
		err = k8s.UpdatePassword(ns, name, username, password, server)
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
