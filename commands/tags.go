package commands

import (
    "errors"
    "github.com/common-nighthawk/go-figure"
    "github.com/fatih/color"
    "github.com/olekukonko/tablewriter"
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
        return errors.New(color.RedString("github token not empty"))
    }

    common.SpinnerStart("doing...")
    newClient(token)
    directoryContent := fetchRepositoryContent()
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
