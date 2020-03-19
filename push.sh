docker build -t nabsul/k8s-awsreg-renew .
docker tag nabsul/k8s-awsreg-renew:latest nabsul/k8s-awsreg-renew:$1
docker push nabsul/k8s-awsreg-renew:$1
