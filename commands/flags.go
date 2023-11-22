package commands

import (
    "github.com/urfave/cli/v2"
    "os/exec"
    "strings"
)

func Flags() []cli.Flag {
    return []cli.Flag{
        &cli.StringFlag{
            Name:     "token",
            Aliases:  []string{"t"},
            Required: false,
            Value:    "",
        },
    }
}

func fetchToken() string {
    cmd := exec.Command("sh", "-c", "git config --list | grep 'search.workflow.token'")
    stdout, err2 := cmd.Output()
    if err2 != nil {
        return ""
    }

    return strings.ReplaceAll(strings.TrimSpace(string(stdout)), "search.workflow.token=", "")
}
