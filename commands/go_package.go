package commands

import (
    "bufio"
    "fmt"
    "github.com/common-nighthawk/go-figure"
    "github.com/pibigstar/termcolor"
    "github.com/urfave/cli/v2"
    "github.com/xiaoxuan6/rsearch/common"
    "io"
    "regexp"
    "strings"
)

var GoPackageCommand = &cli.Command{
    Name:        common.GoCommandName,
    Usage:       common.GoCommandUsage,
    Description: figure.NewFigure("rsearch "+common.GoCommandName, "", true).String() + common.GoCommandUsage,
    Action:      Exec,
}

func Exec(c *cli.Context) error {
    b, err2 := common.Get("https://github.moeyy.xyz/https://raw.githubusercontent.com/xiaoxuan6/go-package-example/main/README.md")

    if err2 != nil {
        return err2
    }

    common.SpinnerStart("sync go doing...")
    var modelsSlice []common.Model
    br := bufio.NewReader(strings.NewReader(string(b)))
    for {
        a, _, errs := br.ReadLine()
        if errs == io.EOF {
            break
        }

        re := regexpContent(string(a))
        if len(re) < 3 {
            continue
        }

        modelsSlice = append(modelsSlice, common.Model{
            Title: re[3],
            Tag:   common.GoTagName,
            Url:   re[2],
        })
    }

    common.SpinnerStop()
    if len(modelsSlice) > 0 {
        _ = common.DeleteByTag(common.GoTagName)
    }

    err2 = common.CreateInBatches(modelsSlice)
    if err2 != nil {
        fmt.Println(termcolor.FgRed("数据插入失败：" + err2.Error()))
        return nil
    }

    fmt.Print(termcolor.FgGreen("sync successfully"))
    return nil
}

func regexpContent(val string) []string {
    re := regexp.MustCompile(`\|(.*?)\|(.*?)\|(.*?)\|`)
    matchers := re.FindStringSubmatch(val)
    if len(matchers) < 1 {
        return nil
    }
    return matchers
}
