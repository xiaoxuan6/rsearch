package commands

import (
    "errors"
    "github.com/olekukonko/tablewriter"
    "github.com/urfave/cli/v2"
    "os"
    "rsearch/common"
    "strings"
)

func Runs(ctx *cli.Context) error {
    token := ctx.String("token")
    if len(token) < 1 {
        token = fetchToken()
    }

    if token == "" {
        return errors.New("github token not empty")
    }

    client := fetchClient(token)
    directoryContent := fetchRepositoryContent(client)

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
