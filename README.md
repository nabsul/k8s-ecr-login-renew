# Renew Kubernetes Docker secrets for AWS ECR 

AWS Elastic Container Registry (ECR) provides a cost-effective private registry for your Docker containers. 
However, ECR Docker credentials 
[expire every 12 hours](https://docs.aws.amazon.com/cli/latest/reference/ecr/get-login.html).

To work around this, I created this small tool to automatically refresh the secret in Kubernetes.
It deploys as a cron job and ensures that your Kubernetes cluster
will always be able to pull Docker images from ECR.

## Running the Example

The following sections describe step-by-step how to set the cron job up and test it.
You should be able to use the `example/service-account.yml` and `example/deploy.yml` files as-is for testing purposes,
but you'll need to fill in you registry's URL before using `example/pod.yml`.

### Create an ECR Instance

I'm not going to describe this in too much details because
there is [plenty of documentation](https://docs.aws.amazon.com/AmazonECR/latest/userguide/what-is-ecr.html) 
that describes how to do this.
But here are a few pointers:

- Create an AWS ECR instance
- Create a repository inside that instance

### Push a Test Image to 

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
kubectl create job --from=cronjob/ecr-renew-demo ecr-renew-demo-manual-1
```

You can view the status and logs of the job with the following commands:

```shell script
kubectl describe job ecr-renew-demo-manual-1
kubectl logs job/ecr-renew-demo-manual-1
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
If it is not provided, the it will fall back to the `default` namespace.

# Automated Testing

Since this tool is mostly "glue" between AWS and Kuberenetes,
I've decided that unit tests are not so useful.
Instead, the tests in this code are designed to run against a real Kubernetes cluster.
It will auto-detect the cluster to use, 
and will refuse to run if the namespaces it uses already exist 
(to avoid accidentally overriding real configurations).

The only prerequisite to running these tests is to have the required AWS secret created.
This can be done with the following command:

```shell script
kubectl create secret -n tests-ecr-renew generic ecr-renew-cred-demo \
  --from-literal=REGION=[AWS_REGION] \
  --from-literal=ID=[AWS_KEY_ID] \
  --from-literal=SECRET=[AWS_SECRET]
```

You can then run the tests by typing:

```shell script
go test -v
```
