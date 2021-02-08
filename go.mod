module github.com/nabsul/k8s-ecr-login-renew

go 1.14

require (
	github.com/aws/aws-sdk-go v1.34.28
	github.com/hashicorp/vault v1.6.1
	github.com/hashicorp/vault/api v1.0.5-0.20201001211907-38d91b749c77
	github.com/imdario/mergo v0.3.8 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v0.18.2
	sigs.k8s.io/structured-merge-diff v0.0.0-20190525122527-15d366b2352e // indirect
)
