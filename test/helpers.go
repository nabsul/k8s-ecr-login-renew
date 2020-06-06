package test

import (
	"github.com/nabsul/k8s-ecr-login-renew/src/k8s"
	v1 "k8s.io/api/core/v1"
	"strings"
	"testing"
)

func CheckNamespaces(t *testing.T) {
	namespaces, err := k8s.GetAllNamespaces()
	if nil != err {
		t.Error(err)
		return
	}

	if 0 == len(namespaces) {
		t.Error("no namespaces returned")
	}

	t.Logf("Namespaces found: %s", strings.Join(namespaces, ", "))
}

func CreateNamespace(name string) (*v1.Namespace, error) {
	c, err := k8s.GetClient()
	if nil != err {
		return nil, err
	}

	ns := &v1.Namespace{}
	ns.Name = name
	return c.CoreV1().Namespaces().Create(ns)
}
