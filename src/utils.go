package src

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func MakePath(appName, path string) string {
	if strings.HasPrefix(path, "/") {
		path = strings.TrimLeft(path, "/")
	}

	if strings.HasSuffix(appName, "/") {
		appName = strings.TrimRight(appName, "/")
	}

	res := strings.Join([]string{appName, path}, "/")
	if strings.HasPrefix(res, "/") {
		return res
	}

	return "/" + res
}

func MakeParameter(path string, nextToken *string) *ssm.GetParametersByPathInput {
	return &ssm.GetParametersByPathInput{
		Path:           aws.String(path + "/"),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
		NextToken:      nextToken,
	}
}

func IsSecured(t string) bool {
	return strings.ToUpper(t) == "SECURESTRING"
}

func GetSecuredTypeString(secured bool) string {
	if secured {
		return "SecureString"
	}

	return "String"
}
