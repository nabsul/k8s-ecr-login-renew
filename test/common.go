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

func printError(t *testing.T, err error) {
	t.Errorf("%+v", err)
}
