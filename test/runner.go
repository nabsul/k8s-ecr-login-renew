package test

import (
	"errors"
	"fmt"
	"github.com/nabsul/k8s-ecr-login-renew/src/k8s"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"testing"
	"time"
)

type config struct {
	t                 *testing.T
	createdNamespaces []string
	successNamespaces []string
	failNamespaces    []string
	canGetNamespaces  bool
	targetNamespace   string
}

func runTest(cfg config) {
	cleanup(cfg.createdNamespaces)
	time.Sleep(10 * time.Second)
	//cfg.t.Cleanup(func(){cleanup(cfg.createdNamespaces)})

	t := cfg.t
	c, err := k8s.GetClient()
	if err != nil {
		printError(t, err)
		return
	}

	err = checkNamespaces(c, cfg.createdNamespaces)
	if err != nil {
		printError(t, err)
		return
	}

	t.Log("Creating namespaces")
	for _, ns := range cfg.createdNamespaces {
		_, err := createNamespace(c, ns)
		if err != nil {
			printError(t, err)
			return
		}
	}

	t.Log("Creating service and permissions")
	err = createServiceAccount(c, cfg.successNamespaces, cfg.canGetNamespaces)
	if err != nil {
		printError(t, err)
		return
	}

	awsParams, err := getAwsParams(c)
	awsId, ok1 := awsParams["ID"]
	awsSecret, ok2 := awsParams["SECRET"]
	awsRegion, ok3 := awsParams["REGION"]
	if !(ok1 && ok2 && ok3) {
		values := fmt.Sprintf("[%s], [%s], [%s]", awsId, awsSecret, awsRegion)
		err = errors.New(fmt.Sprintf("one or more AWS Param not configured: %s", values))
		printError(t, err)
		return
	}

	t.Log("Creating cron job")
	err = initCronJob(c, cfg.targetNamespace, awsRegion, awsId, awsSecret)
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
		printError(t, errors.New(fmt.Sprintf("no AWS success message found")))
	}

	expectedUpdates := len(cfg.successNamespaces) + len(cfg.failNamespaces)
	actualUpdates := strings.Count(logs, "Updating secret in namespace")
	if actualUpdates != expectedUpdates {
		msg := fmt.Sprintf("unexpected number of updates %d != %d", expectedUpdates, actualUpdates)
		printError(t, errors.New(msg))
	}

	for _, ns := range cfg.successNamespaces {
		t.Logf("Checking for success: %s", ns)
		msg := fmt.Sprintf("Updating secret in namespace [%s]... success", ns)
		if !strings.Contains(logs, msg) {
			msg = fmt.Sprintf("no success message found for namespace [%s]", ns)
			printError(t, errors.New(msg))
		}
	}

	for _, ns := range cfg.failNamespaces {
		t.Logf("Checking for failure: %s", ns)
		msg := fmt.Sprintf("Updating secret in namespace [%s]... failed:", ns)
		if !strings.Contains(logs, msg) {
			msg = fmt.Sprintf("no fail message found for namespace [%s]", ns)
			printError(t, errors.New(msg))
		}
	}

	if t.Failed() {
		printError(t, errors.New(logs))
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
			fmt.Printf("Failed to cleanup namespace [%s]: [%s}\n", ns, err)
		}
	}
}
