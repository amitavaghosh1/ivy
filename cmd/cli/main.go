package main

import (
	"ivy/src"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const VERSION = "1.0"

func main() {
	rc := src.ReaderCommand{}
	wc := src.NewWriterCommand(src.YamlParser{})
	dc := src.DiffCommand{}

	app := &cli.App{
		Name:  "ivy",
		Usage: "ssm parameter handling utility",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    src.Region,
				Aliases: []string{src.RegionShort},
				Value:   "us-west-2",
				EnvVars: []string{"AWS_REGION"},
			},
			&cli.StringFlag{
				Name:    src.Version,
				Aliases: []string{src.VersionShort},
				Value:   VERSION,
			},
			&cli.StringFlag{
				Name:    src.AwsProfile,
				Aliases: []string{src.AwsProfileShort},
				EnvVars: []string{"AWS_PROFILE"},
				Value:   "default",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "read",
				Usage: rc.Help(),
				Flags: []cli.Flag{
					&cli.StringFlag{Name: src.Service, Aliases: []string{src.ServiceShort}, Required: true},
					&cli.StringSliceFlag{Name: src.Environment, Aliases: []string{src.EnvironmentShort}, Value: cli.NewStringSlice("dev")},
					&cli.StringFlag{Name: src.Format, Aliases: []string{src.FormatShort}, Value: "yaml"},
					&cli.StringFlag{Name: src.Region, Aliases: []string{src.RegionShort}, Value: "us-west-2"},
					&cli.BoolFlag{Name: src.DecryptSecured, Aliases: []string{src.DecryptSecuredShort}, Value: false},
					&cli.StringFlag{Name: src.AwsProfile, Aliases: []string{src.AwsProfileShort}, EnvVars: []string{"AWS_PROFILE"}, Value: "default"},
				},
				Action: rc.Run,
			},
			{
				Name:  "write",
				Usage: wc.Help(),
				Flags: []cli.Flag{
					&cli.StringFlag{Name: src.File, Aliases: []string{src.FileShort}, Required: true},
					&cli.BoolFlag{Name: src.Overwrite, Aliases: []string{src.OverwriteShort}, Value: false},
					&cli.StringSliceFlag{Name: src.Environment, Aliases: []string{src.EnvironmentShort}},
					&cli.StringFlag{Name: src.AwsProfile, Aliases: []string{src.AwsProfileShort}, EnvVars: []string{"AWS_PROFILE"}, Value: "default"},
				},
				Action: wc.Run,
			},
			{
				Name:  "diff",
				Usage: dc.Help(),
				Flags: []cli.Flag{
					&cli.StringFlag{Name: src.Service, Aliases: []string{src.ServiceShort}, Required: true},
					&cli.StringFlag{Name: src.DiffEnv, Aliases: []string{src.DiffEnvShort}, Required: true},
					&cli.StringFlag{Name: src.Environment, Aliases: []string{src.EnvironmentShort}, Value: "dev"},
					&cli.StringFlag{Name: src.Region, Aliases: []string{src.RegionShort}, Value: "us-west-2"},
					&cli.StringFlag{Name: src.AwsProfile, Aliases: []string{src.AwsProfileShort}, EnvVars: []string{"AWS_PROFILE"}, Value: "default"},
				},
				Action: dc.Run,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
