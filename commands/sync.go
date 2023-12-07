package commands

import (
    "context"
    "errors"
    "fmt"
    "github.com/common-nighthawk/go-figure"
    "github.com/google/go-github/v48/github"
    "github.com/pibigstar/termcolor"
    "github.com/sirupsen/logrus"
    "github.com/urfave/cli/v2"
    "github.com/xiaoxuan6/rsearch/common"
    "golang.org/x/oauth2"
    "regexp"
    "strings"
    "sync"
)

var (
    wg          sync.WaitGroup
    client      *github.Client
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
    newClient(token)
    directoryContent := fetchRepositoryContent()
    for _, val := range directoryContent {
        filename := val.GetName()
        if strings.HasSuffix(filename, ".md") {
            wg.Add(1)
            go fetchUrlContent(c, filename, &wg)
        }
    }

    wg.Wait()
    common.SpinnerStop()

    fmt.Print(termcolor.FgGreen("sync successfully"))
    return nil
}

func newClient(token string) {
    oauth := oauth2.NewClient(c, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
    client = github.NewClient(oauth)
}

func fetchRepositoryContent() []*github.RepositoryContent {
    _, directoryContent, _, _ := client.Repositories.GetContents(c, common.Owner, common.Repo, "", &github.RepositoryContentGetOptions{})
    return directoryContent
}

func fetchUrlContent(ctx context.Context, filename string, wg *sync.WaitGroup) {
    defer wg.Done()

    RepositoryContent, _, _, err2 := client.Repositories.GetContents(ctx, common.Owner, common.Repo, filename, &github.RepositoryContentGetOptions{})
    if err2 != nil {
        logrus.Error(err2.Error())
        return
    }

    content, err3 := RepositoryContent.GetContent()
    if err3 != nil {
        logrus.Error(err3.Error())
        return
    }

    fetchFileContent([]byte(content), strings.ReplaceAll(filename, ".md", ""))
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
