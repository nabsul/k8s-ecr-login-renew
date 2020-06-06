package aws

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"strings"
)

func GetUserAndPass() (username, password, server string, err error) {
	svc := ecr.New(session.Must(session.NewSession()))
	token, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if nil != err {
		return "", "", "", err
	}

	auth := token.AuthorizationData[0]

	decode, err := base64.StdEncoding.DecodeString(*auth.AuthorizationToken)
	if nil != err {
		return "", "", "", err
	}

	parts := strings.Split(string(decode), ":")
	return parts[0], parts[1], *auth.ProxyEndpoint, nil
}
