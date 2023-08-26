package main

import (
    "errors"
    "github.com/common-nighthawk/go-figure"
    "github.com/sirupsen/logrus"
    "github.com/urfave/cli/v2"
    "os"
    "path/filepath"
    "rsearch/commands"
    "rsearch/common"
    "strconv"
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
        param := os.Args[1]
        if param != common.CommandName && param != common.GoCommandName && param != "count" && param != "clear" && param != "resync" {
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
                    // rsearch sync --path="xxx"
                    &cli.StringFlag{
                        Name:     "RepositoryPath",
                        Aliases:  []string{"path"},
                        Required: false,
                        Value:    common.RepositoryPathDefault,
                        EnvVars:  []string{"REPOSITORY_PATH"},
                        Action: func(context *cli.Context, s string) error {
                            if _, err := os.Stat(s); os.IsNotExist(err) {
                                return errors.New("给定参数 `" + s + "` 不存在")
                            }
                            return nil
                        },
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
                Name:        "resync",
                Usage:       "重新同步所有数据",
                Description: figure.NewFigure("rsearch resync", "", true).String(),
                Action: func(context *cli.Context) error {
                    _ = commands.Run(context)
                    _ = commands.Exec(context)
                    return nil
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        logrus.Error(err.Error())
    }
}
