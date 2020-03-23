docker build -t nabsul/k8s-ecr-login-renew .
docker tag nabsul/k8s-ecr-login-renew:latest nabsul/k8s-ecr-login-renew:$1
docker tag nabsul/k8s-ecr-login-renew:latest nabsul/k8s-ecr-login-renew:latest
docker push nabsul/k8s-ecr-login-renew:$1
docker push nabsul/k8s-ecr-login-renew:latest
