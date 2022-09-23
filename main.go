package main

import (
	"os"
	"strings"

	"github.com/nabsul/k8s-ecr-login-renew/src/aws"
	"github.com/nabsul/k8s-ecr-login-renew/src/k8s"
	"github.com/nabsul/k8s-ecr-login-renew/src/utils"
	"go.uber.org/zap"
)

const (
	envVarAwsSecret       = "DOCKER_SECRET_NAME"
	envVarTargetNamespace = "TARGET_NAMESPACE"
	envVarRegistries      = "DOCKER_REGISTRIES"
)

var logger *zap.SugaredLogger = nil
var err error = nil

func init() {
	logger, err = utils.GetLogger()
	if err != nil {
		panic(err)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	logger.Info("starting k8s-ecr-loging-renew...")

	name := os.Getenv(envVarAwsSecret)
	if name == "" {
		logger.Fatalf("environment variable %s is required.", envVarAwsSecret)
	}

	logger.Info("fetching auth data from Amazon Web Services...")
	credentials, err := aws.GetDockerCredentials()
	checkErr(err)
	logger.Info("fetched docker credentials.")

	servers := getServerList(credentials.Server)
	logger.Info("docker registries found:", zap.Any("registries", strings.Join(servers, ",")))

	namespaces, err := k8s.GetNamespaces(os.Getenv(envVarTargetNamespace))
	checkErr(err)
	logger.Info("updating kubernetes secret in specified namespaces...", zap.String("secret", name), zap.Int("namespaces", len(namespaces)))

	failedNamespaces := []string{}
	for _, namespace := range namespaces {
		logger.Info("updating secret in namespace... ", zap.String("secret", name), zap.String("namespace", namespace))
		err = k8s.UpdatePassword(namespace, name, credentials.Username, credentials.Password, servers)
		if nil != err {
			logger.Error("failed to update secret.", zap.String("secret", name), zap.String("namespace", namespace), zap.Error(err))
			failedNamespaces = append(failedNamespaces, namespace)
		} else {
			logger.Info("successfully updated secret", zap.String("secret", name), zap.String("namespace", namespace))
		}
	}

	if len(failedNamespaces) > 0 {
		logger.Fatalf("failed to create or update one or more docker login registry secret in mentioned namespaces", zap.Any("namespaces", failedNamespaces))
	}

	logger.Info("job completed.")
}

func getServerList(defaultServer string) []string {
	addedServersSetting := os.Getenv(envVarRegistries)

	if addedServersSetting == "" {
		return []string{defaultServer}
	}

	return strings.Split(addedServersSetting, ",")
}
