package src

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

type DiffCommand struct{}

func (d DiffCommand) Help() string {
	return strings.TrimSpace(`
	Diff ssm params in two envs
	Usage:
	ivy diff
		-s param_name [csv-backend]
		-e env [dev/test source env to diff]
		-d env [which env to diff with]
		-r region [defaults to us-west-2] (optional)
	`)
}

func (d DiffCommand) Run(c *cli.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	envParams, err := d.getParamsForEnv(ctx, c.String(Environment), c.String(Service), c.String(Region))
	if err != nil {
		return err
	}

	otherParams, err := d.getParamsForEnv(ctx, c.String(DiffEnv), c.String(Service), c.String(Region))
	if err != nil {
		return err
	}

	var b []byte

	difflogs := Diff(envParams, otherParams)

	b, err = json.MarshalIndent(difflogs, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}

func (d DiffCommand) getParamsForEnv(ctx context.Context, env, service, region string) ([]ConfigParameter, error) {
	ac := AwsConfigReader{
		AppName: service,
		Region:  region,
		Path:    env,
	}

	return ac.Read(ctx, true)
}

type Changelog struct {
	Action string
	ConfigParameter
}

func Diff(lparams, rparams []ConfigParameter) []Changelog {
	changes := []Changelog{}

	lmap, rmap := ToMap(lparams), ToMap(rparams)

	for key, value := range lmap {
		config, ok := rmap[key]
		if !ok {
			changes = append(changes, Changelog{Action: "delete", ConfigParameter: value})
		} else if value.Changed(config) {
			changes = append(changes, Changelog{Action: "update", ConfigParameter: value})
		}
	}

	return changes
}

func ToMap(params []ConfigParameter) map[string]ConfigParameter {
	lmap := map[string]ConfigParameter{}

	for i := range params {
		param := params[i]

		lmap[param.Name] = param
	}

	return lmap
}
