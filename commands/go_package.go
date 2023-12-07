package commands

import (
    "bufio"
    "errors"
    "fmt"
    "github.com/avast/retry-go"
    "github.com/common-nighthawk/go-figure"
    "github.com/fatih/color"
    "github.com/sirupsen/logrus"
    "github.com/urfave/cli/v2"
    "github.com/xiaoxuan6/rsearch/common"
    "io"
    "io/ioutil"
    "net/http"
    "regexp"
    "strings"
    "time"
)

var GoPackageCommand = &cli.Command{
    Name:        common.GoCommandName,
    Usage:       common.GoCommandUsage,
    Description: figure.NewFigure("rsearch "+common.GoCommandName, "", true).String() + common.GoCommandUsage,
    Action:      Exec,
}

func Exec(c *cli.Context) error {
    b, err2 := fileGetContent()
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
        fmt.Println(color.RedString("数据插入失败：" + err2.Error()))
        return nil
    }

    fmt.Print(color.GreenString("sync successfully"))
    return nil
}

func fileGetContent() (b []byte, err error) {
    err = retry.Do(
        func() error {
            client := http.Client{
                Timeout: 5 * time.Second,
            }

            response, err1 := client.Get(common.GoPackageRepository)
            if err1 != nil {
                return errors.New("请求错误：" + err1.Error())
            }

            defer func() {
                _ = response.Body.Close()
            }()

            b, err1 = ioutil.ReadAll(response.Body)
            if err1 != nil {
                return errors.New("获取内容失败：" + err1.Error())
            }

            return nil
        },
        retry.Attempts(3),
        retry.OnRetry(func(n uint, err error) {
            logrus.Info(fmt.Sprintf("请求失败，第 %d 次重试", n+1))
        }),
        retry.LastErrorOnly(true),
    )

    if err != nil {
        return b, err
    }

    return b, nil
}

func regexpContent(val string) []string {
    re := regexp.MustCompile(`\|(.*?)\|(.*?)\|(.*?)\|`)
    matchers := re.FindStringSubmatch(val)
    if len(matchers) < 1 {
        return nil
    }
    return matchers
}
