package main

import (
	"fmt"
	"github.com/nabsul/k8s-ecr-login-renew/src/k8s"
	"github.com/nabsul/k8s-ecr-login-renew/test"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func Test_BasicFunction(t *testing.T) {
	namespaces := []string{test.ConstSvcNamespace}
	///initTest(t, namespaces)

	err := runTest(t, namespaces, namespaces, false)
	if err != nil {
		t.Error(err)
	}
}


func runTest(t *testing.T, namespaces, allowedNamespaces []string, canGetNamespaces bool) error {
	t.Log("Creating namespaces")
	for _, ns := range namespaces {
		_, err := test.CreateNamespace(ns)
		if err != nil {
			return err
		}
	}

	t.Log("Creating service and permissions")
	err := test.CreateServiceAccount(allowedNamespaces, canGetNamespaces)
	if err != nil {
		return err
	}

	err = test.CopyAwsSecret()
	if err != nil {
		return err
	}

	t.Log("Creating cron job")
	err = test.CreateCronJob()
	if nil != err {
		return err
	}

	t.Log("Running the job")
	err = test.RunCronJob()
	if nil != err {
		return err
	}

	return nil
}

func initTest(t *testing.T, namespaces []string) {
	t.Cleanup(func(){cleanup(namespaces)})
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
