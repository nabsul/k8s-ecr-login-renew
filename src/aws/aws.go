package aws

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"strings"
)

type EcrCredentials struct {
	Username, Password, Server string
}

func GetUserAndPass() ([]EcrCredentials, error) {
	svc := ecr.New(session.Must(session.NewSession()))
	token, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if nil != err {
		return nil, err
	}

	result := make([]EcrCredentials, len(token.AuthorizationData))
	for _, auth := range token.AuthorizationData {
		decode, err := base64.StdEncoding.DecodeString(*auth.AuthorizationToken)
		if nil != err {
			return nil, err
		}

		parts := strings.Split(string(decode), ":")
		cred := EcrCredentials{
			Username: parts[0],
			Password: parts[1],
			Server: *auth.ProxyEndpoint,
		}

		result = append(result, cred)
	}

	return result, nil
}
