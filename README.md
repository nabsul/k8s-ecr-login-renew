# Renew Kubernetes Docker secrets for AWS ECR 

AWS Elastic Container Registry (ECR) provides a cost-effective private registry for your Docker containers. 
However, ECR Docker credentials 
[expire every 12 hours](https://docs.aws.amazon.com/cli/latest/reference/ecr/get-login.html).

To work around this, I created this small tool to automatically refresh the secret in Kubernetes.
It deploys as a cron job and ensures that your Kubernetes cluster
will always be able to pull Docker images from ECR.

## Quick Start

Prerequisite: AWS IAM credentials with permission to read ECR data.

Installation with Helm:

```sh
helm repo add nabsul https://nabsul.github.io/helm
helm install k8s-ecr-login-renew nabsul/k8s-ecr-login-renew --set awsRegion=[REGION],awsAccessKeyId=[ACCESS_KEY_ID],awsSecretAccessKey=[SECRET_KEY]
```

## Docker Images

The tool is built for and supports the following Architectures:
- `linux/amd64`
- `linux/arm64`
- `linux/arm/v7`

If there is an architecture that isn't supported you can request it [here](https://github.com/nabsul/k8s-ecr-login-renew/issues).

The Docker image for running this tool in Kubernetes is published here: https://hub.docker.com/r/nabsul/k8s-ecr-login-renew

Note1: Although a `latest` tag is currently being published, I highly recommend using a specific version.
With the `latest` you run the risk of using an outdated version of the tool, or getting upgraded to a newer version before you're ready.
I will eventually deprecate the `latest` tag. UPDATE: It happened sort of by accident, but the `latest` tag is now gone and won't be coming back.

## Environment Variables

The tool is configured using the following environment variables:

- AWS_ACCESS_KEY_ID (required): AWS access key used to create the Docker credentials.
- AWS_SECRET_ACCESS_KEY (required): AWS secret needed to fetch Docker credentials from AWS.
- AWS_REGION (required): The AWS region where your ECR instance is created.
- DOCKER_SECRET_NAME (required): The name of the Kubernetes secret where the Docker credentials are stored.
- TARGET_NAMESPACE (optional): Comma, semicolon or newline separated list of namespaces. 
  A Docker secret is created in each of these. 
  If this environment variable is not set, a value of `default` is assumed.
- DOCKER_REGISTRIES (optional): Comma-separated list of registry URL. 
  If none is provided, the default URL returned from AWS is used.
  - Example: `DOCKER_REGISTRIES=https://321321.dkr.ecr.us-west-2.amazonaws.com,https://123123.dkr.ecr.us-east-2.amazonaws.com`

## Prerequisites

The following sections describe step-by-step how to set up the prerequisites needed to deploy this tool.

### Create an ECR Instance

I'm not going to describe this in too much details because there is 
[plenty of documentation](https://docs.aws.amazon.com/AmazonECR/latest/userguide/what-is-ecr.html) 
that describes how to do this. Here are a few pointers:

- Create an AWS ECR instance
- Create a repository inside that instance

### Push a Test Image to ECR 

If you are not already using ECR and pushing images to it,
you'll need to create  and upload an test image to ECR.
Theres's plenty of good documentation out there, but basically:

- Install the AWS CLI tool and Docker on your machine 
- Log into the registry: `$(aws2 ecr get-login --no-include-email --region [AWS_REGION])`
- Build your image: `docker build -t [ECR_URL]:latest .`
- Push your docker image to the registry: `docker push [ECR_URL]:latest`

### Setup AWS Permissions

In AWS IAM, create a user with readonly permissions to your registry.
You should give this account the minimal amount of permissions necessary.
Ideally you should set up a policy that:

- Gives read-only access
- Is restricted to the specific instance(s) of ECR that you want to access

I find IAM to be rather tricky, but here are the steps that I followed:

- Select "Add User", select the "Programmatic Access" option
- (Optionally) Create a group for that user
- Authorize the group or user to pull images from ECR. Either
  a. Use the existing AWS policy "AmazonEC2ContainerRegistryReadOnly"
  b. or create a policy for that group or user with the following configuration:
    - Service: Elastic Container Registry
    - Access Level: List & Read
    - Resources: Select the specific ECR instance that you'll be using
  
Once that's created, you'll want to copy and save the Access Key ID and Secret Access Key of the user for the next step.
I recommend storing these secrets in some kind of secret store, such as: 
[Doppler](https://www.doppler.com/),
[Azure Key Vault](https://azure.microsoft.com/en-us/services/key-vault),
[AWS Secret Store](https://azure.microsoft.com/en-us/services/key-vault),
[1Password](https://1password.com/),
[LastPass](https://www.lastpass.com/)

### Deploy AWS Credentials to Kubernetes as a secret

Note: If you want to use Helm to create this secret automatically, you can skip this section.

You will need to create a secret in Kubernetes with the IAM user's credentials.
The secret can be created from the command line using `kubectl` as follows:

```shell script
kubectl create secret -n ns-ecr-renew-demo generic ecr-renew-cred-demo \
  --from-literal=aws-access-key-id=[AWS_KEY_ID] \
  --from-literal=aws-secret-access-key=[AWS_SECRET]
```

## Deploy to Kubernetes

There are two ways to deploy this tool, and you only need to use one of them:

- Helm chart
- Plain YAML files

### Deploy to Kubernetes with Helm 3

Add the repository

```
helm repo add nabsul https://nabsul.github.io/helm
```

Deploy to your Kubernetes cluster with:

```sh
awsRegion=[REGION],awsAccessKeyId=[ACCESS_KEY_ID],awsSecretAccessKey=[SECRET_KEY]
```

Note: If you have already created a secret with your IAM credentials, you only need to provide a region parameter to Helm.

You can uninstall the tool with:

```sh
helm uninstall k8s-ecr-login-renew
```

### Deploy to Kubernetes with plain YAML

If you don't want to use Helm to manage installing this tool, you can use [`deploy.yaml`](https://github.com/nabsul/k8s-ecr-login-renew/blob/main/deploy.yaml) and `kubectl apply`.
Note that this file is generated from the Helm template by running `helm template .\chart --set forHelm=false,awsRegion=us-west-2 > deploy.yaml`.
You will likely need to review and edit this yaml file to your needs, and then you can deploy with:

```sh
kubectl apply -f deploy.yaml
```

You can also uninstall from your Kubernetes cluster with:

```sh
kubectl delete -f deploy.yaml
```

## Test the Cron Job

To check if the cron job is correctly configured, you can wait for the job to be run.
However, you can also manually trigger a job with the following command:

```sh
kubectl create job --from=cronjob/k8s-ecr-login-renew-cron k8s-ecr-login-renew-cron-manual-1
```

You can view the status and logs of the job with the following commands:

```sh
kubectl describe job k8s-ecr-login-renew-cron-manual-1
kubectl logs job/k8s-ecr-login-renew-cron-manual-1
```

### Deploying ECR Images

You should now be able to deploy them to a pod.
Note that you will need to specify the Docker secret in your Pod definition by adding a `imagePullSecrets` 
field pointing to the created Docker secret (named `k8s-ecr-login-renew-docker-secret` by default).

You can find more information about this here: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/

### Running in a namespace other than default namespace

The example configuration runs in a namespace called `ns-ecr-renew-demo`.
This is configured using the `TARGET_NAMESPACE` environment variable.
If it is not provided, it will fall back to the `default` namespace.

### Multiple Namespace Support

The `TARGET_NAMESPACE` environment variable can also be used to indicate multiple namespaces.
This can be done using comma-separated names as wildcards.
The two wildcard characters are `?` (match a single character) and `*` (match zero or more characters).
The following values are valid for the `TARGET_NAMESPACE` environment variable:

- `namespace1`
- `namespace1,namespace2`
- `*`
- `namespaece1,prefix-*,*-suffix`
- `namespace1,namespace?0`

> Note: If you decide to use wildcards,
> the tool must be granted permission to list all namespaces in your cluster.

# Automated Testing

Since this tool is mostly "glue" between AWS and Kuberenetes,
I've decided that unit tests are not so useful.
Instead, the tests here are designed to run against a real Kubernetes cluster.
It will auto-detect the cluster to use, 
and will refuse to run if the namespaces it uses already exist 
(to avoid accidentally overriding real configurations).

> note: I have only tried running these unit tests on my local Windows machine with  Docker Desktop.

Running these tests locally has the following prerequisites:
 
- Build an image of the tool:  `docker build -t test-ecr-renew .`
- Create a secret with the needed AWS parameters:

```shell script
kubectl create secret generic test-ecr-renew-aws --from-literal=REGION=[AWS_REGION] --from-literal=ID=[AWS_ID] --from-literal=SECRET=[AWS_SECRET] --from-literal=IMAGE=test-ecr-renew
```

You can then run the tests by typing:

```shell script
go test -v ./test/...
```

## Test Limitation

One of the biggest things I am currently not testing are permissions.
This is mostly due to laziness on my part:
I couldn't figure out how to get Docker Desktop to enforce RBAC permissions. 

## Cross-Architecture Build

```shell
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t [image_name] .
```
