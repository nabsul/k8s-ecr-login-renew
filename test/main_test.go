package test

import (
	"strings"
	"testing"
)

func Test_NoTargetNamespace(t *testing.T) {
	runTest(config{
		t:                 t,
		successNamespaces: []string{"default"},
		targetNamespace:   "",
	})
}

func Test_SingleNamespace(t *testing.T) {
	runTest(config{
		t:                 t,
		successNamespaces: []string{ConstSvcNamespace},
		targetNamespace:   ConstSvcNamespace,
	})
}

func Test_SingleNamespace2(t *testing.T) {
	runTest(config{
		t:                 t,
		successNamespaces: []string{spaces["22"]},
		targetNamespace:   spaces["22"],
	})
}

func Test_TwoNamespaces(t *testing.T) {
	namespaces := []string{spaces["3"], spaces["12"]}
	runTest(config{
		t:                 t,
		successNamespaces: namespaces,
		targetNamespace:   strings.Join(namespaces, ","),
	})
}

func Test_WithStar(t *testing.T) {
	namespaces := allNamespaces()
	runTest(config{
		t:                 t,
		successNamespaces: namespaces,
		targetNamespace:   "test-ecr-renew-*",
	})
}

func Test_WithQuestionMark(t *testing.T) {
	runTest(config{
		t:                 t,
		successNamespaces: []string{spaces["12"], spaces["22"]},
		targetNamespace:   "test-ecr-renew-ns?2",
	})
}

func Test_WithStarAndQuestionMark(t *testing.T) {
	runTest(config{
		t:                 t,
		successNamespaces: []string{spaces["12"], spaces["22"]},
		targetNamespace:   "test-ecr-*-ns?2",
	})
}

func Test_OverlappingResults1(t *testing.T) {
	runTest(config{
		t:                 t,
		successNamespaces: allNamespaces(),
		targetNamespace:   "test-ecr-renew-ns13,test-ecr-renew-ns?2,test-ecr-renew-*",
	})
}

func Test_OverlappingResults2(t *testing.T) {
	runTest(config{
		t:                 t,
		successNamespaces: allNamespaces(),
		targetNamespace:   "test-ecr-renew-ns?2,test-ecr-renew-*,test-ecr-renew-ns13",
	})
}

func Test_OverlappingResults3(t *testing.T) {
	runTest(config{
		t:                 t,
		successNamespaces: allNamespaces(),
		targetNamespace:   "test-ecr-renew-*,test-ecr-renew-ns13,test-ecr-renew-ns?2",
	})
}
