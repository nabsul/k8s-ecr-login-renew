package aws

import (
	"encoding/base64"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// GetUserAndPass : function to retrieve the login details to AWS ECR
// Params :
// awsAccessID : AWS access key ID,
// awsSecret : AWS secret access key ,
// region : AWS Region where ECR is setup
// NOTE: if all params are "", then the secrets are taken from AWS environment variables
func GetUserAndPass(awsAccessID, awsSecret, region string) (username, password, server string, err error) {

	var sess *session.Session
	if awsAccessID == "" || awsSecret == "" || region == "" {
		sess, err = session.NewSession()
		if nil != err {
			return "", "", "", err
		}
	} else {
		sess, err = session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(awsAccessID, awsSecret, ""),
		})
		if nil != err {
			return "", "", "", err
		}
	}
	svc := ecr.New(sess)
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
