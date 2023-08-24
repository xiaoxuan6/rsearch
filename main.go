package main

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"rsearch/command"
	"rsearch/common"
	"strconv"
)

func main() {

	common.InitDb(common.SqlitePath)
	defer func() {
		common.CloseDb()
	}()

	if os.Args[1] == "count" {
		num := common.Count()
		logrus.Info("数据库总条数为：" + strconv.Itoa(int(num)))
		os.Exit(0)
	}

	if os.Args[1] != common.CommandName {
		if _, err := os.Stat(common.SqlitePath); os.IsNotExist(err) {
			logrus.Error("sqlite db not exist, 请执行 `rsearch sync` 在重试")
			os.Exit(0)
		}

		command.Search(os.Args[1])
		os.Exit(0)
	}

	app := cli.App{
		Name:        "rsearch",
		Usage:       "rsearch",
		Description: figure.NewFigure("rsearch", "", true).String(),
		Commands: []*cli.Command{
			{
				Name:    common.CommandName,
				Aliases: []string{"s 同步远程数据保存到本地 sqlite 数据库"},
				Description: figure.NewFigure("rsearch sync", "", true).String() +
					"同步远程数据保存到本地 sqlite 数据库",
				Action: command.Run,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Error(err.Error())
	}
}
