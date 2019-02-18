package main

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"strings"
)

func getUserAndPass() (username, password, server string) {
	svc := ecr.New(session.Must(session.NewSession()))
	token, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	checkErr(err)

	auth := token.AuthorizationData[0]

	decode, err := base64.StdEncoding.DecodeString(*auth.AuthorizationToken)
	checkErr(err)

	parts := strings.Split(string(decode), ":")
	return parts[0], parts[1], *auth.ProxyEndpoint
}
