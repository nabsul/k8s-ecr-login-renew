package main

import (
	"fmt"
	"github.com/nabsul/k8s-ecr-login-renew/src/k8s"
	"github.com/nabsul/k8s-ecr-login-renew/test"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestFirst(t *testing.T) {
	initTest(t)
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
	initTest(t)
	err := test.CreateCronJob("default", "test-job")
	if nil != err {
		t.Error(err)
		return
	}
}

func initTest(t *testing.T) {
	t.Cleanup(cleanup)
}

func cleanup() {
	c, err := k8s.GetClient()
	if nil != err {
		return
	}

	c.CoreV1().Namespaces().Delete("ns-ecr-demo-1", &v1.DeleteOptions{})
	c.BatchV1beta1().CronJobs("default").Delete("test-job", &v1.DeleteOptions{})
}
