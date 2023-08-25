package main

import (
    "errors"
    "github.com/common-nighthawk/go-figure"
    "github.com/sirupsen/logrus"
    "github.com/urfave/cli/v2"
    "os"
    "path/filepath"
    "rsearch/command"
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

        switch os.Args[1] {
        case "count":
            num := common.Count()
            logrus.Info("数据库总条数为：" + strconv.Itoa(int(num)))
            os.Exit(0)
        case "clear":
            _ = common.Clear()
            logrus.Info("数据清空成功")
            os.Exit(0)
        }

        if os.Args[1] != common.CommandName {
            command.Search(os.Args[1])
            os.Exit(0)
        }

    }

    app := cli.App{
        Name:        "rsearch",
        Usage:       "rsearch",
        Description: figure.NewFigure("rsearch", "", true).String(),
        Commands: []*cli.Command{
            {
                Name:  common.CommandName,
                Usage: "同步远程数据保存到本地 sqlite 数据库",
                Description: figure.NewFigure("rsearch sync", "", true).String() +
                    "同步远程数据保存到本地 sqlite 数据库",
                Action: command.Run,
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
            },
            {
                Name:  "count",
                Usage: "查询数据总条数",
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        logrus.Error(err.Error())
    }
}
