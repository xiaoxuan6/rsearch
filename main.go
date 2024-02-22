package main

import (
    context2 "context"
    "fmt"
    "github.com/common-nighthawk/go-figure"
    "github.com/mitchellh/go-homedir"
    "github.com/pibigstar/termcolor"
    "github.com/sirupsen/logrus"
    "github.com/urfave/cli/v2"
    "github.com/xiaoxuan6/rsearch/commands"
    "github.com/xiaoxuan6/rsearch/common"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

var version string

func main() {
    dir, err := homedir.Dir()
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(-1)
    }

    basePath := filepath.Join(dir, "/.rsearch")
    common.InitDb(basePath + "/" + common.SqlitePath)
    defer common.CloseDb()

    if len(os.Args) > 1 {
        search()
    }

    app := cli.App{
        Name:        "rsearch",
        Usage:       "rsearch",
        Description: figure.NewFigure("rsearch", "", true).String(),
        Commands: []*cli.Command{
            commands.SyncCommand,
            commands.TagCommand,
            commands.GoPackageCommand,
            {
                Name:  "clear",
                Usage: "清空所有数据",
                Action: func(context *cli.Context) error {
                    _ = common.Flush()
                    fmt.Print(termcolor.FgGreen("数据清空成功！"))
                    return nil
                },
            },
            {
                Name:  "count",
                Usage: "查询数据总条数",
                Action: func(context *cli.Context) error {
                    num := common.Count()
                    fmt.Print("数据库总条数为：" + strconv.Itoa(int(num)))
                    return nil
                },
            },
            {
                Name:        "version",
                Usage:       "查看版本号",
                Aliases:     []string{"v"},
                Description: figure.NewFigure("rsearch version", "", true).String(),
                Action: func(context *cli.Context) error {

                    if len(version) < 1 {
                        common.NewClient(common.GetToken(""))

                        tags, _, errs := common.Client.Repositories.ListTags(context2.Background(), common.Owner, "rsearch", nil)
                        if errs != nil {
                            return nil
                        }

                        version = *tags[0].Name
                    }

                    fmt.Println("rsearch version: " + termcolor.FgGreen(version))
                    return nil
                },
            },
        },
    }

    if err = app.Run(os.Args); err != nil {
        logrus.Error(err.Error())
    }
}

func search() {
    commandNames := []string{common.CommandName, common.GoCommandName, "count", "clear", "tags", "--help", "-h", "v", "version"}

    target := false
    param := os.Args[1]
    for _, val := range commandNames {
        if strings.Compare(param, val) == 0 {
            target = true
            continue
        }
    }

    if target == false {
        tag := ""
        if len(os.Args) == 3 {
            tag = os.Args[2]
        }
        commands.Search(param, tag)
        os.Exit(0)
    }
}
