package test

import (
	"github.com/nabsul/k8s-ecr-login-renew/src/k8s"
	v13 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateCronJob(namespace, name string) error {
	c, err := k8s.GetClient()
	if nil != err {
		return err
	}

	job := getCronJob(namespace, name)
	_, err = c.BatchV1beta1().CronJobs(namespace).Create(&job)
	return err
}

func getCronJob(namespace, name string) v1beta1.CronJob {
	return v1beta1.CronJob{
		TypeMeta: v12.TypeMeta{
			Kind:       "CronJob",
			APIVersion: "batch/v1beta1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "ecr-renew-test",
			},
		},
		Spec: v1beta1.CronJobSpec{
			Schedule: "0 0 1 1 1", // set to a value in the past so it never triggers
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: v13.JobSpec{
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							RestartPolicy:      "OnFailure",
							ServiceAccountName: "svc-ecr-renew-demo",
							Containers:         []v1.Container{getContainer()},
						},
					},
				},
			},
		},
	}
}

func getContainer() v1.Container {
	return v1.Container{
		Name:                     "ecr-renew",
		Image:                    "nabsul/k8s-ecr-login-renew:latest",
		Env:                      getCronJobEnvVars(),
	}
}

func getCronJobEnvVars() []v1.EnvVar {
	return []v1.EnvVar{
		createEnvVar("DOCKER_SECRET_NAME", "ecr-docker-login-demo"),
		createEnvVar("TARGET_NAMESPACE", "ns-ecr-renew-demo"),
		createSecretEnvVar("AWS_REGION", "ecr-renew-cred-demo", "REGION"),
		createSecretEnvVar("AWS_REGION", "ecr-renew-cred-demo", "ID"),
		createSecretEnvVar("AWS_REGION", "ecr-renew-cred-demo", "SECRET"),
	}
}

func createEnvVar(name, value string) v1.EnvVar {
	return v1.EnvVar{
		Name:  name,
		Value: value,
	}
}

func createSecretEnvVar(envName, secretName, secretKey string) v1.EnvVar {
	return v1.EnvVar{
		Name: envName,
		ValueFrom: &v1.EnvVarSource{
			SecretKeyRef: &v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: secretName,
				},
				Key: secretKey,
			},
		},
	}
}