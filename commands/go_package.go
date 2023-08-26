package commands

import (
    "errors"
    "github.com/sirupsen/logrus"
    "github.com/urfave/cli/v2"
    "io/ioutil"
    "net/http"
    "regexp"
    "rsearch/common"
    "strings"
)

func Exec(c *cli.Context) error {
    _ = common.DeleteByTag(common.GoTagName)

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

    err2 = common.CreateInBatches(modelsSlice)
    if err2 != nil {
        logrus.Error("数据插入失败：" + err2.Error())
    }

    return nil
}

func fileGetContent() (b []byte, err error) {
    response, err := http.Get(common.GoPackageRepository)
    if err != nil {
        return b, errors.New("请求错误：" + err.Error())
    }

    defer response.Body.Close()

    b, err = ioutil.ReadAll(response.Body)
    if err != nil {
        return b, errors.New("获取内容失败：" + err.Error())
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
