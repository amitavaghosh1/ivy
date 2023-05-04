package src

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

type ReaderCommand struct{}

func (r ReaderCommand) Help() string {
	return strings.TrimSpace(`
	Get the ssm param configuration.
	Usage:
	ivy read
		-s param_name [notification_service]
		-e env [dev/test/staging/whateverr you have access]
		-t [json/yaml/.env, defaults to yaml] (optional)
		-r region [defaults to us-west-2] (optional)
		-p profile [aws profile, defaults to default]
	`)
}

func (r ReaderCommand) Run(c *cli.Context) error {
	res := map[string]map[string][]ConfigParameter{}

	envs := c.StringSlice(Environment)
	ac := AwsConfigReader{
		AppName: c.String(Service),
		Region:  c.String(Region),
		Profile: c.String(AwsProfile),
	}
	log.Println("using profile", ac.Profile)

	decrypted := c.Bool(DecryptSecured)

	res[ac.AppName] = map[string][]ConfigParameter{}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	for _, env := range envs {
		ac.Path = env

		params, err := ac.Read(ctx, decrypted)
		if err != nil {
			return err
		}

		res[ac.AppName][env] = params
	}

	format := c.String(Format)

	switch format {
	case "env":
		HandleEnv(res[ac.AppName], envs[0])
	case "json":
		HandleJSON(res)
	default:
		HandleYAML(res)
	}

	// fmt.Println(string(b))
	return nil
}
