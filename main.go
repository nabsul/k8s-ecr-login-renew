package main

import (
	"fmt"
	"os"
	"time"
)

const awsEmailEnv = "KUBE_AWSREG_EMAIL"
const awsSecretName = "KUBE_AWSREG_SECRET_NAME"

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
	email := getEnv(awsEmailEnv)

	fmt.Print("Fetching auth data from AWS... ")
	username, password, server := getUserAndPass()
	fmt.Println("Success.")

	fmt.Printf("Updating kubernetes secret [%s]... ", name)
	updatePassword(name, username, password, email, server)
	fmt.Println("Success.")
}
