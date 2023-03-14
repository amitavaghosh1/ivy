package src

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type EnvConfig map[string][]ConfigParameter

type WriterCommand struct {
	parser Parser
}

func NewWriterCommand(parser Parser) WriterCommand {
	return WriterCommand{parser: parser}
}

func (r WriterCommand) Help() string {
	return strings.TrimSpace(`
	Update the ssm param configuration.
	Usage:
	ivy write
		-f yaml_file_name [./somefile.yaml]
		-e [dev/test/staging/whatever you have access] (optional)
		-t [json/yaml, defaults to yaml] (optional)
	`)
}

func (r WriterCommand) Run(c *cli.Context) error {
	configMap, err := r.parser.Parse(c.String(File))
	if err != nil {
		return err
	}

	keys := reflect.ValueOf(configMap).MapKeys()
	if len(keys) == 0 {
		return errors.New("invalid file")
	}

	service := keys[0].String()

	envMap := map[string]struct{}{}
	for _, em := range c.StringSlice(Environment) {
		envMap[em] = struct{}{}
	}

	awsConfigs := map[string][]ConfigParameter{}
	isEnvSpecified := len(envMap) > 0

	// fmt.Println("envmap ", envMap, len(envMap))

	for env, values := range configMap[service] {
		if _, ok := envMap[env]; !ok && isEnvSpecified {
			continue
		}

		configs := []ConfigParameter{}
		for _, value := range values {
			configs = append(configs, ConfigParameter{
				Name:  MakePath(service, MakePath(env, value.Name)),
				Type:  GetSecuredTypeString(value.Secured),
				Value: value.Value,
			})
		}
		awsConfigs[env] = configs
	}

	ac := AwsConfigReader{
		Region:  c.String(Region),
		AppName: service,
	}

	// logrus.Debugf("awsconfigs %+v\n", awsConfigs)

	var wg sync.WaitGroup
	wg.Add(len(awsConfigs))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	for k, cfgs := range awsConfigs {
		logrus.Printf("Updating for %s env\n", k)

		go func(ctx context.Context, ac AwsConfigReader, cfgs []ConfigParameter, overwrite bool) {
			defer wg.Done()

			logrus.Infof("%+v", cfgs)

			if err := ac.Write(ctx, cfgs, overwrite); err != nil {
				logrus.Error(err)
			}
		}(ctx, ac, cfgs, c.Bool(Overwrite))
	}
	wg.Wait()

	logrus.Println("success")
	return nil
}
