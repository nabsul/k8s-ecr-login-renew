package test

import (
	"context"
	"errors"
	"fmt"
	"time"

	batchV1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func initCronJob(c *kubernetes.Clientset, targetNamespace, awsRegion, awsId, awsSecret string) error {
	ctx := context.Background()
	job := createCronJob(targetNamespace, awsRegion, awsId, awsSecret)
	_, err := c.BatchV1beta1().CronJobs(ConstSvcNamespace).Create(ctx, &job, metaV1.CreateOptions{})
	return err
}

func runCronJob(c *kubernetes.Clientset) (string, error) {
	ctx := context.Background()
	getOpt := metaV1.GetOptions{}
	cron, err := c.BatchV1().CronJobs(ConstSvcNamespace).Get(ctx, ConstCronJobName, getOpt)
	if err != nil {
		return "", err
	}

	job := createJob(*cron)
	run, err := c.BatchV1().Jobs(ConstSvcNamespace).Create(ctx, &job, metaV1.CreateOptions{})
	if err != nil {
		return "", err
	}

	checkCount := 0
	for run.Status.CompletionTime == nil {
		time.Sleep(5 * time.Second)
		run, err = c.BatchV1().Jobs(ConstSvcNamespace).Get(ctx, job.Name, getOpt)
		if err != nil {
			return "", err
		}

		checkCount += 1
		if checkCount >= 10 {
			return "", errors.New("job ran for too long")
		}
	}

	listOpt := metaV1.ListOptions{LabelSelector: "job-name=test-ecr-renew-job"}
	list, err := c.CoreV1().Pods(ConstSvcNamespace).List(ctx, listOpt)
	if err != nil {
		return "", err
	}

	if len(list.Items) != 1 {
		return "", errors.New(fmt.Sprint("Unexpected number of pods returned from job: %i", len(list.Items)))
	}

	pod := list.Items[0]

	req := c.CoreV1().Pods(ConstSvcNamespace).GetLogs(pod.Name, &coreV1.PodLogOptions{})
	res := req.Do(ctx)
	bytes, err := res.Raw()
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func createJob(cron batchV1.CronJob) batchV1.Job {
	return batchV1.Job{
		TypeMeta: metaV1.TypeMeta{
			Kind: "Job",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: ConstSvcNamespace,
			Name:      ConstJobName,
		},
		Spec: cron.Spec.JobTemplate.Spec,
	}
}

func createCronJob(targetNamespace, awsRegion, awsId, awsSecret string) v1beta1.CronJob {
	one := int32(1)
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
			ConcurrencyPolicy: "Forbid",
			Schedule:          "0 0 1 1 1", // set to a value in the past so it never triggers
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: batchV1.JobSpec{
					Parallelism:  &one,
					Completions:  &one,
					BackoffLimit: &one,
					Template: coreV1.PodTemplateSpec{
						Spec: getPodSpec(targetNamespace, awsRegion, awsId, awsSecret),
					},
				},
			},
		},
	}
}

func getPodSpec(targetNamespace, awsRegion, awsId, awsSecret string) coreV1.PodSpec {
	return coreV1.PodSpec{
		RestartPolicy:      "Never",
		ServiceAccountName: ConstSvcName,
		Containers: []coreV1.Container{
			{
				Name:            "ecr-renew",
				Image:           "test-ecr-renew",
				ImagePullPolicy: "IfNotPresent",
				Env: []coreV1.EnvVar{
					{Name: "DOCKER_SECRET_NAME", Value: ConstDockerSecretName},
					{Name: "TARGET_NAMESPACE", Value: targetNamespace},
					{Name: "AWS_REGION", Value: awsRegion},
					{Name: "AWS_ACCESS_KEY_ID", Value: awsId},
					{Name: "AWS_SECRET_ACCESS_KEY", Value: awsSecret},
				},
			},
		},
	}
}
