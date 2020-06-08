package test

import (
	coreV1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func createServiceAccount(c *kubernetes.Clientset, allowedNamespaces []string, canGetNamespaces bool) error {
	account := &coreV1.ServiceAccount{
		TypeMeta: metaV1.TypeMeta{
			Kind: "ServiceAccount",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: ConstSvcNamespace,
			Name: ConstSvcName,
		},
	}

	_, err := c.CoreV1().ServiceAccounts(ConstSvcNamespace).Create(account)
	if err != nil {
		return err
	}

	/*
	for _, ns := range allowedNamespaces {
		role := createRole(ns)
		_, err = c.RbacV1().Roles(ns).Create(role)
		if err != nil {
			return err
		}

		binding := createRoleBinding(role, svc)
		_, err = c.RbacV1().RoleBindings(ns).Create(binding)
		if err != nil {
			return err
		}
	}

	if canGetNamespaces {
		role := createNamespaceRole()
		_, err = c.RbacV1().ClusterRoles().Create(&role)
		if err != nil {
			return err
		}

		binding := createNamespaceRoleBinding()
		_, err = c.RbacV1().ClusterRoleBindings().Create(&binding)
		if err != nil {
			return err
		}
	}
	 */

	return nil
}

func createNamespaceRole() rbacV1.ClusterRole {
	return rbacV1.ClusterRole{
		TypeMeta: metaV1.TypeMeta{
			Kind: "ClusterRole",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: ConstSvcNamespace,
			Name: ConstNamespaceRoleName,
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs: []string{"list"},
			},
		},
	}
}

func createNamespaceRoleBinding() rbacV1.ClusterRoleBinding {
	return rbacV1.ClusterRoleBinding{
		TypeMeta: metaV1.TypeMeta{
			Kind: "ClusterRoleBinding",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: ConstSvcNamespace,
			Name:      ConstNamespaceRoleBinding,
		},
		RoleRef: rbacV1.RoleRef{
			Kind: "ClusterRoleBinding",
			Name: ConstNamespaceRoleName,
		},
		Subjects: []rbacV1.Subject{
			{
				Namespace: ConstSvcNamespace,
				Name: ConstSvcName,
			},
		},
	}
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
			Namespace: role.Namespace,
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
