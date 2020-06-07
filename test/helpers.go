package test

import (
	"errors"
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func checkNamespaces(c *kubernetes.Clientset, testNamespaces []string) error {
	namespaces, err := c.CoreV1().Namespaces().List(metaV1.ListOptions{})
	if nil != err {
		return err
	}

	if 0 == len(namespaces.Items) {
		return errors.New("no namespaces returned")
	}

	set := map[string]bool{}
	for _, ns := range testNamespaces {
		set[ns] = true
	}

	for _, ns := range namespaces.Items {
		_, ok := set[ns.Name]
		if ok {
			return errors.New(fmt.Sprintf("Namespace already exists: [%s]", ns.Name))
		}
	}

	return nil
}

func createNamespace(c *kubernetes.Clientset, name string) (*coreV1.Namespace, error) {
	ns := &coreV1.Namespace{}
	ns.Name = name
	return c.CoreV1().Namespaces().Create(ns)
}
