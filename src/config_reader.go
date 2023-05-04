package src

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/sirupsen/logrus"
)

type ConfigParameter struct {
	Name    string `json:"name" yaml:"name"`
	Secured bool   `json:"secured" yaml:"secured"`
	Type    string `json:"-" yaml:"-"`
	Value   string `json:"value" yaml:"value"`
}

func (c ConfigParameter) Changed(d ConfigParameter) (res bool) {
	return c.Name != d.Name ||
		c.Secured != d.Secured ||
		c.Type != d.Type
}

type AwsConfigReader struct {
	Region  string
	AppName string
	Path    string
	Profile string
}

func NewReader(region, appName, path, profile string) *AwsConfigReader {
	if strings.TrimSpace(region) == "" {
		region = DefaultRegion
	}

	if strings.TrimSpace(profile) == "" {
		profile = "default"
	}

	return &AwsConfigReader{
		Region:  region,
		AppName: appName,
		Path:    path,
		Profile: profile,
	}
}

func (ar *AwsConfigReader) GetPath() string {
	return MakePath(ar.AppName, ar.Path)
}

func (ar *AwsConfigReader) Read(ctx context.Context, decrypt bool) ([]ConfigParameter, error) {
	var pp []ConfigParameter

	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(ar.Region)},
		SharedConfigState: session.SharedConfigEnable,
		Profile:           ar.Profile,
	})

	if err != nil {
		return pp, err
	}

	p := ar.GetPath()
	ssmSession := ssm.New(sess, aws.NewConfig().WithRegion(ar.Region).WithMaxRetries(10))

	var nextToken *string
	var params []*ssm.Parameter
	var depth = 0
	const max_depth = 200

	for {
		depth += 1
		out, err := ssmSession.GetParametersByPathWithContext(ctx, MakeParameter(p, nextToken))
		if err != nil {
			return pp, err
		}

		if out.NextToken != nil {
			nextToken = aws.String(*out.NextToken)
		}

		params = append(params, out.Parameters...)
		if out.NextToken == nil || depth >= max_depth {
			break
		}
	}

	for _, res := range params {
		param := ConfigParameter{
			Name:    strings.TrimPrefix(*res.Name, p+"/"),
			Secured: IsSecured(*res.Type),
			Value:   *res.Value,
		}

		pp = append(pp, param)
	}

	return pp, nil
}

func (ar *AwsConfigReader) Write(ctx context.Context, params []ConfigParameter, overwrite bool) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(ar.Region)},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}

	ssmSession := ssm.New(sess, aws.NewConfig().WithRegion(ar.Region).WithMaxRetries(10))

	for _, cfg := range params {
		logrus.Debugln(cfg.Name, "=>", cfg.Value)
		if _, err := ssmSession.PutParameterWithContext(ctx, &ssm.PutParameterInput{
			Name:      aws.String(cfg.Name),
			Value:     aws.String(cfg.Value),
			Type:      aws.String(cfg.Type),
			Overwrite: aws.Bool(overwrite),
		}); err != nil {
			return err
		}
	}

	logrus.Println("success")
	return nil
}
