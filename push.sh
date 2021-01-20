docker build -t nabsul/k8s-ecr-login-renew:$1 .
docker build -t nabsul/k8s-ecr-login-renew:arm32v7-$1 -f Dockerfile-arm .
docker push nabsul/k8s-ecr-login-renew:$1
docker push nabsul/k8s-ecr-login-renew:arm32v7-$1
