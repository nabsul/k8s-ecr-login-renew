package main

import (
	"fmt"
	"os"
	"time"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func checkEnv(names []string) {
	for _, name := range names {
		if os.Getenv(name) == "" {
			panic("Environment variable " + name + " is required")
		}
	}
}

func main() {
	fmt.Println("Running at " + time.Now().UTC().String())
	checkEnv([]string{"KUBE_AWSREG_SECRET_NAME", "KUBE_AWSREG_EMAIL"})

	name, email := os.Getenv("KUBE_AWSREG_SECRET_NAME"), os.Getenv("KUBE_AWSREG_EMAIL")

	fmt.Print("Fetching auth data from AWS... ")
	username, password, server := getUserAndPass()
	fmt.Println("Success.")

	fmt.Print("Updating kubernetes... ")
	updatePassword(name, username, password, email, server)
	fmt.Println("Success.")
}
