package test

import (
	"strings"
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

func Test_OneSuccessOneFail(t *testing.T) {
	namespaces := []string{spaces["_"], spaces["1"], spaces["2"]}
	runTest(config{
		t:                 t,
		createdNamespaces: namespaces,
		successNamespaces: []string{spaces["_"], spaces["1"]},
		failNamespaces:    []string{spaces["2"]},
		canGetNamespaces:  false,
		targetNamespace:   strings.Join(namespaces, ","),
	})
}
