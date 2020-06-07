package main

import (
	"fmt"
	"github.com/nabsul/k8s-ecr-login-renew/src/k8s"
	"github.com/nabsul/k8s-ecr-login-renew/test"
	rbacVa "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestFirst(t *testing.T) {
	initTest()
	name := "ns-ecr-demo-1"
	ns, err := test.CreateNamespace(name)
	if nil != err {
		t.Error(err)
		return
	}

	if name != ns.Name {
		t.Error(fmt.Sprintf("names not matching [%s] != [%s]", name, ns.Name))
		return
	}
}

func Test_CreateCronJob(t *testing.T) {
	initTest()
	err := test.CreateCronJob("default", "test-job")
	if nil != err {
		t.Error(err)
		return
	}
}


func runTest(namespaces []string, roles []*rbacVa.RoleRef) error {
	for _, ns := range namespaces {
		_, err := test.CreateNamespace(ns)
		if err != nil {
			return err
		}
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
