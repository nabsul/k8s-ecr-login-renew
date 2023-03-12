package test

import (
	"context"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getAwsParams(c *kubernetes.Clientset) (map[string]string, error) {
	secret, err := c.CoreV1().Secrets("default").Get(context.Background(), ConstAwsSecretName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for k, v := range secret.Data {
		result[k] = string(v)
	}

	return result, nil
}
