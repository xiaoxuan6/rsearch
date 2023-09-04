package commands

import (
    "context"
    "errors"
    "fmt"
    "github.com/google/go-github/v48/github"
    "github.com/sirupsen/logrus"
    "github.com/urfave/cli/v2"
    "golang.org/x/oauth2"
    "regexp"
    "rsearch/common"
    "strings"
    "sync"
)

var wg sync.WaitGroup
var c = context.Background()

func Run(ctx *cli.Context) error {

    token := ctx.String("token")
    if len(token) < 1 {
        token = FetchToken()
    }

    if token == "" {
        return errors.New("github token not empty")
    }

    err = common.Clear()
    if err != nil {
        return errors.New(fmt.Sprintf("清空数据失败：%s", err.Error()))
    }

    client := fetchClient(token)
    directoryContent := fetchRepositoryContent(client)
    for _, val := range directoryContent {
        filename := val.GetName()
        if strings.HasSuffix(filename, ".md") {
            wg.Add(1)
            logrus.Info("正在同步文件：" + filename)
            go fetchUrlContent(c, client, filename, &wg)
        }
    }

    wg.Wait()
    logrus.Info("sync successfully")

    return nil
}

func fetchClient(token string) *github.Client {
    oauth := oauth2.NewClient(c, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
    client := github.NewClient(oauth)

    return client
}

func fetchRepositoryContent(client *github.Client) []*github.RepositoryContent {
    _, directoryContent, _, _ := client.Repositories.GetContents(c, common.Owner, common.Repo, "", &github.RepositoryContentGetOptions{})
    return directoryContent
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
