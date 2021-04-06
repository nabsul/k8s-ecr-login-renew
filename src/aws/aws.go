package aws

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"strings"
)

type DockerCredentials struct {
	Username, Password, Server string
}

func GetDockerCredentials() (*DockerCredentials, error) {
	svc := ecr.New(session.Must(session.NewSession()))
	token, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if nil != err {
		return nil, err
	}

	// We expect the response to always be a single entry
	// info: https://github.com/nabsul/k8s-ecr-login-renew/pull/18
	auth := token.AuthorizationData[0]

	decode, err := base64.StdEncoding.DecodeString(*auth.AuthorizationToken)
	if nil != err {
		return nil, err
	}

	parts := strings.Split(string(decode), ":")
	cred := DockerCredentials{
		Username: parts[0],
		Password: parts[1],
		Server: *auth.ProxyEndpoint,
	}

	return &cred, nil
}
