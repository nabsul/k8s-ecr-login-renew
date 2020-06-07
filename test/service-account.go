package test

import (
	"github.com/nabsul/k8s-ecr-login-renew/src/k8s"
	coreV1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ConstSvcNamespace = "ns-test-ecr-renew-job"
const ConstSvcName = "test-ecr-renew-svc"
const ConstRoleName = "test-ecr-renew-role"
const ConstRoleBindingName = "test-ecr-renew-role-binding"

func CreateServiceAccount(allowedNamespaces []string) error {

	c, err := k8s.GetClient()
	if nil != err {
		return err
	}

	account := &coreV1.ServiceAccount{
		TypeMeta: metaV1.TypeMeta{
			Kind: "ServiceAccount",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: ConstSvcNamespace,
			Name: ConstSvcName,
		},
	}

	svc, err := c.CoreV1().ServiceAccounts(ConstSvcNamespace).Create(account)
	if err != nil {
		return err
	}

	for _, ns := range allowedNamespaces {
		role := createRole(ns)
		_, err = c.RbacV1().Roles(ns).Create(role)
		if err != nil {
			return err
		}

		binding := createRoleBinding(role, svc)
		_, err = c.RbacV1().RoleBindings(ConstSvcNamespace).Create(binding)
		if err != nil {
			return err
		}
	}

	return nil
}

func createRole(ns string) *rbacV1.Role {
	return &rbacV1.Role{
		TypeMeta: metaV1.TypeMeta{
			Kind: "Role",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: ns,
			Name:      ConstRoleName,
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups:     []string{""},
				Resources:     []string{"secrets"},
				ResourceNames: []string{"test-ecr-renew-docker-login"},
				Verbs:         []string{"get", "delete"},
			},
			{
				APIGroups:     []string{""},
				Resources:     []string{"secrets"},
				ResourceNames: []string{"test-ecr-renew-docker-login"},
				Verbs:         []string{"create"},
			},
		},
	}
}

func createRoleBinding(role *rbacV1.Role, service *coreV1.ServiceAccount) *rbacV1.RoleBinding {
	return &rbacV1.RoleBinding{
		TypeMeta:   metaV1.TypeMeta{
			Kind: "RoleBinding",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: role.Name,
			Name: ConstRoleBindingName,
		},
		RoleRef:    rbacV1.RoleRef{
			Kind: role.Kind,
			Name: role.Name,
		},
		Subjects:   []rbacV1.Subject{
			{
				Kind: "ServiceAccount",
				Namespace: service.Namespace,
				Name: service.Name,
			},
		},
	}
}
