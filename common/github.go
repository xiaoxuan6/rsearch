package common

import (
    "context"
    "github.com/gofri/go-github-ratelimit/github_ratelimit"
    "github.com/google/go-github/v48/github"
    "github.com/sirupsen/logrus"
    "golang.org/x/oauth2"
    "strings"
)

var (
    c      = context.Background()
    Client *github.Client
)

func NewClient(token string) {
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(c, ts)
    rateLimiter, errs := github_ratelimit.NewRateLimitWaiterClient(tc.Transport)
    if errs != nil {
        panic(errs.Error())
    }

    Client = github.NewClient(rateLimiter)
}

func FetchRepositoryContent() []*github.RepositoryContent {
    _, directoryContent, _, _ := Client.Repositories.GetContents(c, Owner, Repo, "", &github.RepositoryContentGetOptions{})
    return directoryContent
}

func FetchUrlContent(ctx context.Context, filename string) ([]byte, string, error) {
    RepositoryContent, _, _, err2 := Client.Repositories.GetContents(ctx, Owner, Repo, filename, &github.RepositoryContentGetOptions{})
    if err2 != nil {
        logrus.Error(err2.Error())
        return nil, "", err2
    }

    content, err3 := RepositoryContent.GetContent()
    if err3 != nil {
        logrus.Error(err3.Error())
        return nil, "", err3
    }

    return []byte(content), strings.ReplaceAll(filename, ".md", ""), nil
}
