package test

import (
	"github.com/nabsul/k8s-ecr-login-renew/src/k8s"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ConstSecretName = "test-ecr-renew-aws"

func CopyAwsSecret() error {
	c, err := k8s.GetClient()
	if err != nil {
		return err
	}

	secret, err := c.CoreV1().Secrets("default").Get(ConstSecretName, v1.GetOptions{})
	if err != nil {
		return err
	}
	secret.Namespace = ConstSvcNamespace

	newSecret := v12.Secret{
		TypeMeta: v1.TypeMeta{Kind: "Secret"},
		ObjectMeta: v1.ObjectMeta{Namespace: ConstSvcNamespace, Name: ConstSecretName},
		Data: secret.Data,
	}
	_, err = c.CoreV1().Secrets(ConstSvcNamespace).Create(&newSecret)
	if err != nil {
		return err
	}

	return nil
}
