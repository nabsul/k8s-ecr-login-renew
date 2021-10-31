# Renew Kubernetes Docker secrets for AWS ECR 

AWS Elastic Container Registry (ECR) provides a cost-effective private registry for your Docker containers. 
However, ECR Docker credentials 
[expire every 12 hours](https://docs.aws.amazon.com/cli/latest/reference/ecr/get-login.html).

To work around this, I created this small tool to automatically refresh the secret in Kubernetes.
It deploys as a cron job and ensures that your Kubernetes cluster
will always be able to pull Docker images from ECR.

## Docker Images

The tool is built for and supports the following Architectures:
- `linux/amd64`
- `linux/arm64`
- `linux/arm/v7`

If there is an achitecture that isnt supported you can request it [Here](https://github.com/nabsul/k8s-ecr-login-renew/issues).

The latest image can be pulled by any supported Architecture:
- `nabsul/k8s-ecr-login-renew:latest`

Or by tag:
- `nabsul/k8s-ecr-login-renew:v1.5`

## Environment Variables

The tool is mainly configured through environment variables. These are:

- AWS_ACCESS_KEY_ID (required): AWS access key used to create the Docker credentials.
- AWS_SECRET_ACCESS_KEY (required): AWS secret needed to fetch Docker credentials from AWS.
- AWS_REGION (required): The AWS region where your ECR instance is created.
- DOCKER_SECRET_NAME (required): The name of the Kubernetes secret where the Docker credentials are stored.
- TARGET_NAMESPACE (optional): Comma-separated list of namespaces. 
  A Docker secret is created in each of these. 
  If this environment variable is not set, a value of `default` is assumed.
- DOCKER_REGISTRIES (optional): Comma-separated list of registry URL. 
  If none is provided, the default URL returned from AWS is used.
  - Example: `DOCKER_REGISTRIES=https://321321.dkr.ecr.us-west-2.amazonaws.com,https://123123.dkr.ecr.us-east-2.amazonaws.com`

## Running the Example

The following sections describe step-by-step how to set the cron job up and test it.
You should be able to use the `example/service-account.yml` and `example/deploy.yml` files as-is for testing purposes,
but you'll need to fill in you registry's URL before using `example/pod.yml`.

### Create an ECR Instance

I'm not going to describe this in too much details because
there is 
[plenty of documentation](https://docs.aws.amazon.com/AmazonECR/latest/userguide/what-is-ecr.html) 
that describes how to do this.
Here are a few pointers:

- Create an AWS ECR instance
- Create a repository inside that instance

### Push a Test Image to ECR 

To complete the final steps of these instructions, you'll need to create and upload an image to ECR.
As with the previous section, there's plenty of good documentation out there.
But if you're looking to quickly try things out, I've included a trivial Dockerfile
that builds an nginx server image with no modifications:

- Install the AWS CLI tool and Docker on your machine 
- Log into the registry: `$(aws2 ecr get-login --no-include-email --region [AWS_REGION])`
- Build your image: `docker build -t [ECR_URL]:latest example/.`
- Push your docker image to the registry: `docker push [ECR_URL]:latest`

### Setup AWS Permissions

In AWS IAM, create a user with readonly permissions to your registry.
You should give this account the minimal amount of permissions necessary.
Ideally you should set up a policy that:

- Gives read-only access
- Is restricted to the specific instance(s) of ECR that you want to access

I find IAM to be rather tricky, but here are the steps that I followed:

- Select "Add User", select the "Programmatic Access" option
- Create a group for that user
- Create a policy for that group with the following configuration:
  - Service: Elastic Container Registry
  - Access Level: List & Read
  - Resources: Select the specific ECR instance that you'll be using
  
Once that's created, you'll want to get the access key ID and secret for the next step.

### Deploy AWS Access Keys to Kubernetes

You will then need to create a secret in Kubernetes with the IAM user's credentials.
The secret can be created from the command line using `kubectl` as follows:

```shell script
kubectl create secret -n ns-ecr-renew-demo generic ecr-renew-cred-demo \
  --from-literal=REGION=[AWS_REGION] \
  --from-literal=ID=[AWS_KEY_ID] \
  --from-literal=SECRET=[AWS_SECRET]
```

### Required Kubernetes Service Account

You will need to setup a service account with permissions to create/delete/get the resource docker secret.
Ideally, you should give this service account the minimal amount of permissions needed to do its job.
An example of this minimal permissions setup can be found in `example/service-account.yml`.
You can use this apply this configuration directly as follows:

```shell script
kubectl apply -f example/service-account.yml
```

### Deploy the cron job

You'll need to 
Deploy the cron job using the example yaml file in `example/deploy.yml`:

```shell script
kubectl apply -f example/deploy.yml
```

### Test the Cron Job

The easiest way to test is to manually trigger the cron job from the Kubernetes dashboard.
This should create a job and you can then check the logs for any error messages.
Once the job completes, you should notice that the target docker secrets object was either created or updated.

The job can also be manually triggered with the following command:

```shell script
kubectl create job --from=cronjob/cron-ecr-renew-demo cron-ecr-renew-demo-manual-1
```

You can view the status and logs of the job with the following commands:

```shell script
kubectl describe job cron-ecr-renew-demo-manual-1
kubectl logs job/cron-ecr-renew-demo-manual-1
```

### Deploy the Test Image

If you pushed an image to your ECR registry, you should now be able to deploy that image to a pod.
If you edit `example/pod.ym` and replace `[ECR_URI]` with your registry's URI,
you should now be able to run a pod with this command:

```shell script
kubectl apply -f example/pod.yml
```

Check that the pod is running with the following commands:

```shell script
kubectl exec -it ecr-image-pull-test-demo bash
```

This should log you into the running pod, where you can execute commands such as `ls`, 
`cat /usr/share/nginx/html/index.html` and `exit`.
You can also try running the following:

```shell script
kubectl port-forward ecr-image-pull-test-demo 8080:80
```

And then you can open http://localhost:8080 in your browser to see an nginx default welcome message.

### Clean up after the Demo

After running this demo, you might want to clean up everything.
Since the demo is all in its own namespace, just delete it:

```shell script
kubectl delete namespace ns-ecr-renew-demo
```

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
kubectl create secret test-ecr-renew-aws-settings \
  --from-literal=REGION=[AWS_REGION] \
  --from-literal=ID=[AWS_KEY_ID] \
  --from-literal=SECRET=[AWS_SECRET]
  --from-literal=IMAGE=[TEST_IMAGE]
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
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 --push -t nabsul/k8s-ecr-login-renew:v1.6-rc1 .
```
