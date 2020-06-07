package test

import (
	"strings"
	"testing"
)

func Test_BasicFunction(t *testing.T) {
	RunTest(Config{
		T: t,
		SuccessNamespaces: []string{ConstSvcNamespace},
		FailNamespaces: []string{},
		CanGetNamespaces: false,
		TargetNamespace: ConstSvcNamespace,
	})
}

func Test_OneSuccessOneFail(t *testing.T) {
	namespaces := []string{ConstSvcNamespace, "test-ecr-renew-ns1"}
	RunTest(Config{
		T: t,
		SuccessNamespaces: namespaces,
		FailNamespaces: []string{"test-ecr-renew-ns2"},
		CanGetNamespaces: false,
		TargetNamespace: strings.Join(namespaces, ","),
	})
}
