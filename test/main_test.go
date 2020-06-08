package test

import (
	"testing"
)

var spaces = map[string]string {
	"_": ConstSvcNamespace,
	"1": "test-ecr-renew-ns1",
	"2": "test-ecr-renew-ns2",
	"3": "test-ecr-renew-ns3",
	"11": "test-ecr-renew-ns11",
	"12": "test-ecr-renew-ns12",
	"13": "test-ecr-renew-ns13",
	"21": "test-ecr-renew-ns21",
	"22": "test-ecr-renew-ns22",
	"23": "test-ecr-renew-ns23",
}

func Test_BasicFunction(t *testing.T) {
	runTest(config{
		t:                 t,
		createdNamespaces: []string{ConstSvcNamespace},
		successNamespaces: []string{ConstSvcNamespace},
		failNamespaces:    []string{},
		canGetNamespaces:  false,
		targetNamespace:   ConstSvcNamespace,
	})
}

func Test_DeployToAll(t *testing.T) {
	namespaces := allNamespaces()
	runTest(config{
		t:                 t,
		createdNamespaces: namespaces,
		successNamespaces: namespaces,
		failNamespaces:    []string{},
		canGetNamespaces:  false,
		targetNamespace:   "test-ecr-renew-*",
	})
}

func allNamespaces() []string {
	result := make([]string, 0, len(spaces))
	for _, v := range spaces {
		result = append(result, v)
	}
	return result
}
