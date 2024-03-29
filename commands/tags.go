package commands

import (
    "errors"
    "github.com/common-nighthawk/go-figure"
    "github.com/olekukonko/tablewriter"
    "github.com/pibigstar/termcolor"
    "github.com/urfave/cli/v2"
    "github.com/xiaoxuan6/rsearch/common"
    "os"
    "strings"
)

const tagCommandName = "tags"

var TagCommand = &cli.Command{
    Name:        tagCommandName,
    Usage:       "获取所有的标签",
    Description: figure.NewFigure("rsearch "+tagCommandName, "", true).String(),
    Action:      Runs,
    Flags:       Flags,
}

func Runs(ctx *cli.Context) error {
    token := common.GetToken(ctx.String("token"))
    if token == "" {
        return errors.New(termcolor.FgRed("github token not empty"))
    }

    common.SpinnerStart("doing...")
    common.NewClient(token)
    directoryContent := common.FetchRepositoryContent()
    common.SpinnerStop()

    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"标签"})
    for _, val := range directoryContent {
        if strings.HasSuffix(val.GetName(), ".md") {
            table.Append([]string{strings.ReplaceAll(val.GetName(), ".md", "")})
        }
    }
    table.Append([]string{common.GoTagName})
    table.Render()

    return nil
}
