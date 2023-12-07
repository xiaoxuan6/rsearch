package commands

import (
    "github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
    &cli.StringFlag{
        Name:     "token",
        Aliases:  []string{"t"},
        Required: false,
        Value:    "",
    },
}
