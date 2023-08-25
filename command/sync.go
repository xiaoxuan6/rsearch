package command

import (
    "errors"
    "github.com/sirupsen/logrus"
    "github.com/urfave/cli/v2"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "rsearch/common"
    "strings"
    "sync"
)

func Run(ctx *cli.Context) error {
    common.RepositoryPath = ctx.String("RepositoryPath")
    if _, err := os.Stat(common.RepositoryPath); os.IsNotExist(err) {
        return errors.New("文件夹 `" + common.RepositoryPath + "` 不存在, 请指定参数 `path`")
    }

    files, err := fetchFiles()
    if err != nil {
        return err
    }

    if len(files) < 1 {
        return errors.New("文件夹中不包含 `.md` 文件")
    }

    err = common.Clear()
    if err != nil {
        logrus.Error("清空数据失败：" + err.Error())
    }

    var wg sync.WaitGroup
    for _, file := range files {
        logrus.Info("正在同步文件：" + file)
        //dir, _ := os.Getwd()
        //file := dir + common.RepositoryPath + "/README.md"
        b, err := ioutil.ReadFile(file)
        if err != nil {
            logrus.Printf("文件 %s 读取失败：%s", file, err.Error())
            continue
        }

        wg.Add(1)
        tagName := strings.ToLower(strings.Trim(filepath.Base(file), ".md"))
        go fetchFileContent(&wg, b, tagName)
    }

    wg.Wait()
    logrus.Info("sync successfully")

    return err
}

func fetchFiles() (files []string, err error) {
    path := common.RepositoryPath
    dirFiles, err := ioutil.ReadDir(path)
    if err != nil {
        return files, err
    }

    for _, file := range dirFiles {
        fileName := file.Name()
        if strings.ToLower(filepath.Ext(fileName)) == ".md" {
            files = append(files, path+"/"+fileName)
        }
    }

    return files, err
}

func fetchFileContent(wg *sync.WaitGroup, b []byte, tag string) {
    defer wg.Done()
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
        logrus.Error("数据插入失败：" + err2.Error())
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
