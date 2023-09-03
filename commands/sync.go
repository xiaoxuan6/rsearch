package commands

import (
    "context"
    "errors"
    "fmt"
    "github.com/google/go-github/v48/github"
    "github.com/sirupsen/logrus"
    "github.com/urfave/cli/v2"
    "golang.org/x/oauth2"
    "io/ioutil"
    "regexp"
    "rsearch/common"
    "strings"
    "sync"
)

var wg sync.WaitGroup

func Run(ctx *cli.Context) error {

    token := ctx.String("token")
    if token == "" {
        return errors.New("github token not empty")
    }

    err = common.Clear()
    if err != nil {
        return errors.New(fmt.Sprintf("清空数据失败：%s", err.Error()))
    }

    c := context.Background()
    oauth := oauth2.NewClient(c, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
    client := github.NewClient(oauth)

    b, err2 := ioutil.ReadFile(common.RepositoryFilename)
    if err2 != nil {
        return errors.New(fmt.Sprintf("读取文件失败：%s", err2.Error()))
    }

    content := strings.Split(strings.TrimSpace(string(b)), "\n")
    for _, val := range content {
        wg.Add(1)
        go fetchUrlContent(c, client, val, &wg)
    }

    wg.Wait()
    logrus.Info("sync successfully")

    return nil
}

func fetchUrlContent(ctx context.Context, client *github.Client, filename string, wg *sync.WaitGroup) {
    RepositoryContent, _, _, err2 := client.Repositories.GetContents(ctx, common.Owner, common.Repo, filename, &github.RepositoryContentGetOptions{})
    if err2 != nil {
        logrus.Error(err2.Error())
        wg.Done()
        return
    }

    content, err3 := RepositoryContent.GetContent()
    if err3 != nil {
        wg.Done()
        return
    }

    fetchFileContent(wg, []byte(content), strings.ReplaceAll(filename, ".md", ""))
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
