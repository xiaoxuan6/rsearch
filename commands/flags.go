package commands

import "github.com/urfave/cli/v2"

func Flags() (f []cli.Flag) {
    f = append(f,
        &cli.StringFlag{
            Name:     "token",
            Aliases:  []string{"t"},
            Required: false,
            Value:    "",
        },
    )

    return
}
