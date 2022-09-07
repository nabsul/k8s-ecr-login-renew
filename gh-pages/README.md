# k8s-ecr-login-renew

This tool is designed to run as a cron job in your Kubernetes cluster.
It will periodically connect to AWS, fetch new Docker credentials for your Elastic Container Registries (ECR), and save them in Kubernetes secrets.

This avoids the problem of these credentials expiring every 6 hours.
