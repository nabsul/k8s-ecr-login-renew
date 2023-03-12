package test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nabsul/k8s-ecr-login-renew/src/k8s"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func runTest(cfg config) {
	cfg.t.Cleanup(func() { cleanup(cfg) })

	t := cfg.t
	c, err := k8s.GetClient()
	if err != nil {
		printError(t, err)
		return
	}

	err = checkNamespaces(c, allNamespaces())
	if err != nil {
		printError(t, err)
		return
	}

	t.Log("creating namespaces")
	for _, ns := range allNamespaces() {
		_, err := createNamespace(c, ns)
		if err != nil {
			printError(t, err)
			return
		}
	}

	t.Log("creating service and permissions")
	err = createServiceAccount(c, []string{}, false)
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

	t.Log("creating cron job")
	err = initCronJob(c, cfg.targetNamespace, awsRegion, awsId, awsSecret)
	if nil != err {
		printError(t, err)
		return
	}

	t.Log("running the job")
	logs, err := runCronJob(c)
	if nil != err {
		printError(t, err)
		return
	}

	t.Log("checking job logs")
	if !strings.Contains(logs, "Fetching auth data from AWS... Success.") {
		printError(t, errors.New(fmt.Sprintf("no AWS success message found in:\n%s", logs)))
	}

	expectedUpdates := len(cfg.successNamespaces)
	actualUpdates := strings.Count(logs, "Updating secret in namespace")
	if actualUpdates != expectedUpdates {
		msg := fmt.Sprintf("unexpected number of updates %d != %d", expectedUpdates, actualUpdates)
		printError(t, errors.New(msg))
	}

	for _, ns := range cfg.successNamespaces {
		t.Logf("checking for success: %s", ns)
		msg := fmt.Sprintf("Updating secret in namespace [%s]... success", ns)
		if !strings.Contains(logs, msg) {
			msg = fmt.Sprintf("no success message found for namespace [%s]", ns)
			printError(t, errors.New(msg))
		}
	}

	/* Not currently implemented
	for _, ns := range cfg.failNamespaces {
		t.Logf("checking for failure: %s", ns)
		msg := fmt.Sprintf("Updating secret in namespace [%s]... failed:", ns)
		if !strings.Contains(logs, msg) {
			msg = fmt.Sprintf("no fail message found for namespace [%s]", ns)
			printError(t, errors.New(msg))
		}
	}
	*/

	if t.Failed() {
		printError(t, errors.New(logs))
	}
}

// as long as we only create stuff in namespaces, deleting the namespaces should
func cleanup(cfg config) {
	cfg.t.Log("cleaning up...")

	c, err := k8s.GetClient()
	if nil != err {
		return
	}

	for _, ns := range allNamespaces() {
		err = c.CoreV1().Namespaces().Delete(context.Background(), ns, metaV1.DeleteOptions{})
		if err != nil {
			fmt.Printf("failed to cleanup namespace [%s]: [%s}\n", ns, err)
		}
	}

	cfg.t.Log("giving namespaces time to go away...")
	time.Sleep(10 * time.Second)
}
