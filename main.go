package main

import (
    "github.com/common-nighthawk/go-figure"
    "github.com/olekukonko/tablewriter"
    "github.com/sirupsen/logrus"
    "github.com/urfave/cli/v2"
    "io/ioutil"
    "os"
    "path/filepath"
    "rsearch/commands"
    "rsearch/common"
    "strconv"
    "strings"
)

func main() {
    firstArg := os.Args[0]
    _, file := filepath.Split(firstArg)

    var basePath string
    if file == "main.exe" {
        basePath = "."
    } else {
        ex, _ := os.Executable()
        basePath = filepath.Dir(ex)
    }

    common.InitDb(basePath + "/" + common.SqlitePath)
    defer func() {
        common.CloseDb()
    }()

    if len(os.Args) > 1 {
        commandNames := []string{common.CommandName, common.GoCommandName, "count", "clear", "tags", "--help", "-h"}

        target := false
        param := os.Args[1]
        for _, val := range commandNames {
            if strings.Compare(param, val) == 0 {
                target = true
                continue
            }
        }

        if target == false {
            if param == common.GoTagName {
                commands.TermRenderer()
                os.Exit(0)
            }

            tag := ""
            if len(os.Args) == 3 {
                tag = os.Args[2]
            }
            commands.Search(param, tag)
            os.Exit(0)
        }
    }

    app := cli.App{
        Name:        "rsearch",
        Usage:       "rsearch",
        Description: figure.NewFigure("rsearch", "", true).String(),
        Commands: []*cli.Command{
            {
                Name:        common.CommandName,
                Usage:       common.CommandUsage,
                Description: figure.NewFigure("rsearch sync", "", true).String() + common.CommandUsage,
                Action:      commands.Run,
                Flags: []cli.Flag{
                    &cli.StringFlag{
                        Name:     "token",
                        Aliases:  []string{"t"},
                        Required: false,
                        Value:    "",
                    },
                },
            },
            {
                Name:  "clear",
                Usage: "清空所有数据",
                Action: func(context *cli.Context) error {
                    _ = common.Flush()
                    logrus.Info("数据清空成功")
                    return nil
                },
            },
            {
                Name:  "count",
                Usage: "查询数据总条数",
                Action: func(context *cli.Context) error {
                    num := common.Count()
                    logrus.Info("数据库总条数为：" + strconv.Itoa(int(num)))
                    return nil
                },
            },
            {
                Name:        common.GoCommandName,
                Usage:       common.GoCommandUsage,
                Description: figure.NewFigure("rsearch sync-go", "", true).String() + common.GoCommandUsage,
                Action:      commands.Exec,
            },
            {
                Name:        "tags",
                Usage:       "获取所有的标签",
                Description: figure.NewFigure("rsearch tags", "", true).String(),
                Action: func(context *cli.Context) error {
                    b, _ := ioutil.ReadFile(common.RepositoryFilename)
                    content := strings.Split(strings.TrimSpace(string(b)), "\n")

                    table := tablewriter.NewWriter(os.Stdout)
                    table.SetHeader([]string{"标签"})
                    for _, val := range content {
                        table.Append([]string{strings.ReplaceAll(val, ".md", "")})
                    }
                    table.Render()

                    return nil
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        logrus.Error(err.Error())
    }
}
