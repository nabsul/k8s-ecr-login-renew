package test

import (
	"errors"
	"fmt"
	"github.com/nabsul/k8s-ecr-login-renew/src/k8s"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"testing"
)

type Config struct {
	T                                 *testing.T
	SuccessNamespaces, FailNamespaces []string
	CanGetNamespaces                  bool
	TargetNamespace                   string
}

func RunTest(cfg Config) {
	t := cfg.T
	namespaces := append(cfg.SuccessNamespaces, cfg.FailNamespaces...)

	c, err := k8s.GetClient()
	if err != nil {
		printError(t, err)
		return
	}

	err = checkNamespaces(c, namespaces)
	if err != nil {
		printError(t, err)
		return
	}

	t.Cleanup(func(){cleanup(namespaces)})
	t.Log("Creating namespaces")
	for _, ns := range namespaces {
		_, err := createNamespace(c, ns)
		if err != nil {
			printError(t, err)
			return
		}
	}

	t.Log("Creating service and permissions")
	err = createServiceAccount(c, cfg.SuccessNamespaces, cfg.CanGetNamespaces)
	if err != nil {
		printError(t, err)
		return
	}

	awsParams, err := getAwsParams(c)
	awsId, ok1 := awsParams["ID"]
	awsSecret, ok2 := awsParams["SECRET"]
	awsRegion, ok3 := awsParams["REGION"]
	awsImage, ok4 := awsParams["IMAGE"]
	if !(ok1 && ok2 && ok3 && ok4) {
		values := fmt.Sprintf("[%s], [%s], [%s], [%s]", awsId, awsSecret, awsRegion, awsImage)
		err = errors.New(fmt.Sprintf("one or more AWS Param not configured: %s", values))
		printError(t, err)
		return
	}

	t.Log("Creating cron job")
	err = initCronJob(c, cfg.TargetNamespace, awsRegion, awsId, awsSecret)
	if nil != err {
		printError(t, err)
		return
	}

	t.Log("Running the job")
	logs, err := runCronJob(c)
	if nil != err {
		printError(t, err)
		return
	}

	t.Log("Checking job logs")
	if !strings.Contains(logs, "Fetching auth data from AWS... Success.") {
		printError(t, errors.New(fmt.Sprintf("no AWS success message found in \n%s", logs)))
	}

	for _, ns := range cfg.SuccessNamespaces {
		msg := fmt.Sprintf("Updating secret [%s] in namespace [%s]... success", ConstDockerSecretName, ns)
		if !strings.Contains(logs, msg) {
			msg = fmt.Sprintf("no success message found for namespace [%s] in \n%s", ns, logs)
			printError(t, errors.New(msg))
			return
		}
	}

	for _, ns := range cfg.FailNamespaces {
		msg := fmt.Sprintf("Updating secret [%s] in namespace [%s]... failed:", ConstDockerSecretName, ns)
		if !strings.Contains(logs, msg) {
			msg = fmt.Sprintf("no success message found for namespace [%s] in \n%s", ns, logs)
			printError(t, errors.New(msg))
			return
		}
	}
}

// as long as we only create stuff in namespaces, deleting the namespaces should
func cleanup(namespaces []string) {
	c, err := k8s.GetClient()
	if nil != err {
		return
	}

	for _, ns := range namespaces {
		err = c.CoreV1().Namespaces().Delete(ns, &metaV1.DeleteOptions{})
		if err != nil {
			fmt.Printf("Failed to cleanup namespace [%s]: [%s}", ns, err)
		}
	}
}
