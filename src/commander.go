package src

import "github.com/urfave/cli/v2"

type Commander interface {
	Help() string
	Run(c *cli.Context) error
}
