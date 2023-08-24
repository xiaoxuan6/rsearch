package command

import (
    "github.com/olekukonko/tablewriter"
    "github.com/sirupsen/logrus"
    "os"
    "rsearch/common"
    "strings"
)

var err error
var models []*common.Model

func Search(keyword string) {
    if strings.ToLower(keyword) == "all" {
        models, err = common.All()
    } else {
        models, err = common.Search(keyword)
    }

    if err != nil {
        logrus.Error("fetch data fail err: " + err.Error())
        return
    }

    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"标题", "标签", "地址"})
    for _, val := range models {
        table.Append([]string{val.Title, val.Tag, val.Url})
    }
    table.Render()
}
