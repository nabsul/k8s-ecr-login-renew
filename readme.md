# Renew Kubernetes Docker secrets for AWS ECR 

AWS Elastic Container Registry (ECR) provides a cost-effective private registry for your Docker containers. 
However, ECR Docker credentials 
[expire every 12 hours](https://docs.aws.amazon.com/cli/latest/reference/ecr/get-login.html).

To work around this, I created this small tool to automatically refresh the secret in Kubernetes.
It deploys as a cron job and ensures that your Kubernetes cluster
will always be able to pull Docker images from ECR.

For more information/tutorial, check out this blog post:

https://nabeel.dev/2020/03/22/aws-ecr-renew

## Setup Instructions

The following sections describe step-by-step how to set the cron job up and test it.
You should be able to use the `service-account.yml` and `deploy.yml` files as-is for testing purposes,
but you'll probably want to tweak them for your specific needs before using in production.

I do make the assumption that you've already created an ECR instance and uploaded docker images to it.
There is [plenty of documentation](https://docs.aws.amazon.com/AmazonECR/latest/userguide/what-is-ecr.html) 
that describes how to do that.

### Setup AWS Permissions

In AWS IAM, create a user with readonly permissions to your registry.
You should give this account the minimal amount of permissions necessary.
Ideally you should set up a policy that:

- Gives read-only access
- Is restricted to the specific instance(s) of ECR that you want to access

### Deploy AWS Access Keys to Kubernetes

You will then need to create a secret in Kubernetes with the IAM user's credentials.
The secret can be created from the command line using `kubectl` as follows:

```shell script
kubectl create secret generic ecr-renew-cred --from-literal=ID=[AWS_KEY_ID] --from-literal=SECRET=[AWS_SECRET]
```

### Required Kubernetes Service Account

You will need to setup a service account with permissions to create/delete/get the resource docker secret.
Ideally, you should give this service account the minimal amount of permissions needed to do its job.
An example of this minimal permissions setup can be found in `example/service-account.yml`.
In this file, we only give the service account permission to access a single secret.

### Deploy the cron job

Deploy the cron job using the example yaml file in `example/deploy.yml`:

```shell script
kubectl apply -f example/deploy.yml
```

### Required Environment Variables

The container requires the following environment variables to run:

- `KUBE_AWSREG_SECRET_NAME`: Where to store the docker secret that is created 
- `AWS_REGION`: The region of your AWS ECR instance
- `AWS_ACCESS_KEY_ID`: They AWS Account key ID
- `AWS_SECRET_ACCESS_KEY`: The AWS Account secret key

## Testing the deployment

The easiest way to test is to manually trigger the cron job from the kubernetes dashboard.
This should create a job and you can then check the logs for any error messages.
Once the job completes, you should notice that the target docker secrets object was either created or updated.
