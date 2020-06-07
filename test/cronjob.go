package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nabsul/k8s-ecr-login-renew/src/k8s"
	batchV1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const ConstCronJobName = "test-ecr-renew-cron-job"
const ConstJobName = "test-ecr-renew-job"

func CreateCronJob() error {
	c, err := k8s.GetClient()
	if nil != err {
		return err
	}

	job := createCronJob()
	_, err = c.BatchV1beta1().CronJobs(ConstSvcNamespace).Create(&job)
	return err
}

func RunCronJob() error {
	c, err := k8s.GetClient()
	if nil != err {
		return err
	}

	getOpt := metaV1.GetOptions{}
	cron, err := c.BatchV1beta1().CronJobs(ConstSvcNamespace).Get(ConstCronJobName, getOpt)
	if err != nil {
		return err
	}

	job := createJob(*cron)
	run, err := c.BatchV1().Jobs(ConstSvcNamespace).Create(&job)
	if err != nil {
		return err
	}

	for run.Status.CompletionTime == nil {
		bytes, err := json.MarshalIndent(run, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))

		time.Sleep(5 * time.Second)
		run, err = c.BatchV1().Jobs(ConstSvcNamespace).Get(job.Name, getOpt)
		if err != nil {
			return err
		}
	}
	bytes, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))

	list, err := c.CoreV1().Pods(ConstSvcNamespace).List(metaV1.ListOptions{LabelSelector: "job-name=test-ecr-renew-job"})
	if err != nil {
		return err
	}

	if len(list.Items) != 1 {
		return errors.New(fmt.Sprint("Unexpected number of pods returned: %i", len(list.Items)))
	}

	pod := list.Items[0]

	req := c.CoreV1().Pods(ConstSvcNamespace).GetLogs(pod.Name, &coreV1.PodLogOptions{})
	res := req.Do()
	bytes, err = res.Raw()
	if err != nil {
		return err
	}
	fmt.Printf(string(bytes))

	return nil
}

func createJob(cron v1beta1.CronJob) batchV1.Job {
	return batchV1.Job{
		TypeMeta: metaV1.TypeMeta{
			Kind: "Job",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: ConstSvcNamespace,
			Name: ConstJobName,
		},
		Spec: cron.Spec.JobTemplate.Spec,
	}
}

func createCronJob() v1beta1.CronJob {
	return v1beta1.CronJob{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "CronJob",
			APIVersion: "batch/v1beta1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      ConstCronJobName,
			Namespace: ConstSvcNamespace,
			Labels: map[string]string{
				"app": "test-ecr-renew",
			},
		},
		Spec: v1beta1.CronJobSpec{
			Schedule: "0 0 1 1 1", // set to a value in the past so it never triggers
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: batchV1.JobSpec{
					Template: coreV1.PodTemplateSpec{
						Spec: coreV1.PodSpec{
							RestartPolicy:      "OnFailure",
							ServiceAccountName: ConstSvcName,
							Containers: []coreV1.Container{
								{
									Name:  "ecr-renew",
									Image: "nabsul/k8s-ecr-login-renew:latest",
									Env: []coreV1.EnvVar{
										{Name: "DOCKER_SECRET_NAME", Value: "test-ecr-renew-docker-login"},
										{Name: "TARGET_NAMESPACE", Value: ConstSvcNamespace},
										createSecretEnvVar("AWS_REGION", "REGION"),
										createSecretEnvVar("AWS_ACCESS_KEY_ID", "ID"),
										createSecretEnvVar("AWS_SECRET_ACCESS_KEY", "SECRET"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func createSecretEnvVar(envName, secretKey string) coreV1.EnvVar {
	return coreV1.EnvVar{
		Name: envName,
		ValueFrom: &coreV1.EnvVarSource{
			SecretKeyRef: &coreV1.SecretKeySelector{
				LocalObjectReference: coreV1.LocalObjectReference{
					Name: "test-ecr-renew-aws",
				},
				Key: secretKey,
			},
		},
	}
}