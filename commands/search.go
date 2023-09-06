package commands

import (
    "fmt"
    "github.com/charmbracelet/glamour"
    "github.com/olekukonko/tablewriter"
    "github.com/pibigstar/termcolor"
    "github.com/sirupsen/logrus"
    "os"
    "rsearch/common"
    "strings"
)

var (
    err    error
    models []*common.Model
)

func Search(keyword, tag string) {
    target := false
    if strings.ToLower(keyword) == "all" {
        models, err = common.All()
    } else if tag != "" && keyword == "" {
        models, err = common.FetchModelByTag(tag)
    } else if tag != "" {
        models, err = common.SearchWithTag(keyword, tag)
    } else {
        target = true
        models, err = common.Search(keyword)
    }

    if err != nil {
        logrus.Error("fetch data fail err: " + err.Error())
        return
    }

    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"标题", "标签", "地址"})
    table.SetRowLine(true)
    for _, val := range models {

        var title string
        if target {
            title = replaceKeyword(val.Title, keyword)
        } else {
            title = val.Title
        }

        tag = val.Tag
        if strings.ToLower(keyword) == strings.ToLower(tag) {
            tag = termcolor.BgRed(tag)
        }

        url := val.Url
        if strings.HasPrefix(url, "github") {
            url = "https://" + url
        }

        table.Append([]string{title, tag, termcolor.FgGreen(url)})
    }
    table.Render()
}

func replaceKeyword(title, keyword string) string {
    var keywords []string
    if _, ok := common.FgColor[keyword]; ok {
        keywords = common.FgColor[keyword]
    } else {
        keywords = []string{keyword}
    }

    for _, val := range keywords {
        title = strings.ReplaceAll(title, val, termcolor.BgRed(val))
    }

    return title
}

func TermRenderer() {
    b, err2 := fileGetContent()
    if err2 != nil {
        logrus.Error(err.Error())
        return
    }

    tr, _ := glamour.NewTermRenderer(
        glamour.WithWordWrap(150),
        glamour.WithStylePath("dark"),
    )
    out, _ := tr.Render(string(b))
    fmt.Print(out)
}
