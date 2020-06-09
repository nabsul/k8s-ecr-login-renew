package test

import "testing"

const ConstCronJobName = "test-ecr-renew-cron-job"
const ConstJobName = "test-ecr-renew-job"
const ConstSvcNamespace = "test-ecr-renew-namespace"
const ConstSvcName = "test-ecr-renew-svc"
const ConstRoleName = "test-ecr-renew-role"
const ConstRoleBindingName = "test-ecr-renew-role-binding"
const ConstNamespaceRoleName = "test-ecr-renew-cluster-role"
const ConstNamespaceRoleBinding = "test-ecr-renew-cluster-role-binding"
const ConstAwsSecretName = "test-ecr-renew-aws"
const ConstDockerSecretName = "test-ecr-renew-docker-login"

type config struct {
	t                 *testing.T
	successNamespaces []string
	targetNamespace   string
}

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

func allNamespaces() []string {
	result := make([]string, 0, len(spaces))
	for _, v := range spaces {
		result = append(result, v)
	}
	return result
}

func printError(t *testing.T, err error) {
	t.Errorf("Error: %v", err)
}
