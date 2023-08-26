package commands

import (
    "fmt"
    "github.com/charmbracelet/glamour"
    "github.com/olekukonko/tablewriter"
    "github.com/sirupsen/logrus"
    "io/ioutil"
    "net/http"
    "os"
    "rsearch/common"
    "strings"
)

var err error
var models []*common.Model

func Search(keyword, tag string) {
    if strings.ToLower(keyword) == "all" {
        models, err = common.All()
    } else if tag != "" && keyword == "" {
        models, err = common.FetchModelByTag(tag)
    } else if tag != "" {
        models, err = common.SearchWithTag(keyword, tag)
    } else {
        models, err = common.Search(keyword)
    }

    if err != nil {
        logrus.Error("fetch data fail err: " + err.Error())
        return
    }

    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"标题", "标签", "地址"})
    table.SetRowLine(true)

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

func TermRenderer() {
    response, err := http.Get(common.GoPackageRepository)
    if err != nil {
        logrus.Error("请求错误：" + err.Error())
        return
    }

    defer response.Body.Close()

    b, err := ioutil.ReadAll(response.Body)
    if err != nil {
        logrus.Error("获取内容失败：" + err.Error())
        return
    }

    tr, _ := glamour.NewTermRenderer(
        glamour.WithWordWrap(150),
        glamour.WithStylePath("dark"),
    )
    out, _ := tr.Render(string(b))
    fmt.Print(out)
}
