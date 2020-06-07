package test

import (
	"testing"
)

func Test_BasicFunction(t *testing.T) {
	RunTest(Config{
		T: t,
		SuccessNamespaces: []string{ConstSvcNamespace},
		FailNamespaces: []string{},
		CanGetNamespaces: false,
	})
}
