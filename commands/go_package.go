package commands

import (
    "errors"
    "fmt"
    "github.com/avast/retry-go"
    "github.com/sirupsen/logrus"
    "github.com/urfave/cli/v2"
    "github.com/xiaoxuan6/rsearch/common"
    "io/ioutil"
    "net/http"
    "regexp"
    "strings"
)

func Exec(c *cli.Context) error {
    b, err2 := fileGetContent()
    if err2 != nil {
        return err2
    }

    newContent := string(b)
    replacements := []string{"# Go 开源第三方包收集和使用示例", "|分支名|包名|描述|", "|:---|:---|:---|"}
    for _, replaceOld := range replacements {
        newContent = strings.ReplaceAll(newContent, replaceOld, ``)
    }
    newContent = strings.Trim(newContent, "\n")
    contents := strings.Split(newContent, "\n")
    var modelsSlice []common.Model
    for _, val := range contents {
        regexpStr := regexpContent(val)
        if regexpStr != nil {
            modelsSlice = append(modelsSlice, common.Model{
                Title: regexpStr[3],
                Tag:   common.GoTagName,
                Url:   regexpStr[2],
            })
        }
    }
    
    if len(modelsSlice) > 0 {
        _ = common.DeleteByTag(common.GoTagName)
    }
    err2 = common.CreateInBatches(modelsSlice)
    if err2 != nil {
        logrus.Error("数据插入失败：" + err2.Error())
    }

    return nil
}

func fileGetContent() (b []byte, err error) {
    err = retry.Do(
        func() error {
            response, err1 := http.Get(common.GoPackageRepository)
            if err1 != nil {
                return errors.New("请求错误：" + err1.Error())
            }

            defer response.Body.Close()
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
