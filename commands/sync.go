package commands

import (
    "context"
    "errors"
    "fmt"
    "github.com/common-nighthawk/go-figure"
    "github.com/pibigstar/termcolor"
    "github.com/urfave/cli/v2"
    "github.com/xiaoxuan6/rsearch/common"
    "regexp"
    "strings"
    "sync"
)

var (
    wg          sync.WaitGroup
    c           = context.Background()
    SyncCommand = &cli.Command{
        Name:        common.CommandName,
        Usage:       common.CommandUsage,
        Description: figure.NewFigure("rsearch "+common.CommandName, "", true).String() + common.CommandUsage,
        Action:      Run,
        Flags:       Flags,
    }
)

func Run(ctx *cli.Context) error {
    token := common.GetToken(ctx.String("token"))
    if token == "" {
        return errors.New(termcolor.FgRed("github token not empty"))
    }

    err = common.Clear()
    if err != nil {
        return errors.New(termcolor.FgRed(fmt.Sprintf("清空数据失败：%s", err.Error())))
    }

    common.SpinnerStart("sync doing...")
    common.NewClient(token)
    directoryContent := common.FetchRepositoryContent()
    for _, val := range directoryContent {
        filename := val.GetName()
        if strings.HasSuffix(filename, ".md") {
            wg.Add(1)
            go func(filename string) {
                defer wg.Done()
                body, tag, err1 := common.FetchUrlContent(c, filename)
                if err1 != nil {
                    return
                }

                fetchFileContent(body, tag)
            }(filename)
        }
    }

    wg.Wait()
    common.SpinnerStop()

    fmt.Print(termcolor.FgGreen("sync successfully"))
    return nil
}

func fetchFileContent(b []byte, tag string) {
    strContent := strings.Split(string(b), "\n")
    var ms []common.Model
    for _, value := range strContent {
        url := regexpUrl(value)
        if len(url) < 1 {
            continue
        }

        title := regexpTitle(value)
        ms = append(ms, common.Model{
            Title: title,
            Tag:   tag,
            Url:   url,
        })
    }

    err2 := common.CreateInBatches(ms)
    if err2 != nil {
        fmt.Println(termcolor.FgRed("数据插入失败：" + err2.Error()))
    }
}

func regexpTitle(str string) string {
    re := regexp.MustCompile(`\[(.*?)\]`)
    matches := re.FindStringSubmatch(str)
    if len(matches) > 1 {
        return matches[1]
    }

    return ""
}

func regexpUrl(str string) string {
    re := regexp.MustCompile(`\((.*?)\)`)
    matches := re.FindStringSubmatch(str)
    if len(matches) > 1 {
        return matches[1]
    }
    return ""
}
