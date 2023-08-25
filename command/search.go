package command

import (
    "github.com/charmbracelet/glamour"
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

    tr, _ := glamour.NewTermRenderer(
        glamour.WithAutoStyle(),
        glamour.WithWordWrap(-1),
    )
    for _, val := range models {
        renderUrl, _ := tr.Render(val.Url)
        table.Append([]string{val.Title, val.Tag, renderUrl})
    }
    table.Render()
}
