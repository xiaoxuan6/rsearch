package common

import (
    "bufio"
    "github.com/mitchellh/go-homedir"
    "io"
    "io/ioutil"
    "net/url"
    "os"
    "strings"
)

func GetToken(token string) string {
    if len(token) > 0 {
        return token
    }

    if token = os.Getenv("GITHUB_TOKEN"); len(token) > 0 {
        return token
    }

    gitCredentials, err := homedir.Expand("~/.git-credentials")
    if err != nil {
        return ""
    }

    body, err := ioutil.ReadFile(gitCredentials)
    if err != nil {
        return ""
    }

    r := bufio.NewReader(strings.NewReader(string(body)))
    for {
        line, _, err1 := r.ReadLine()
        if err1 == io.EOF {
            break
        }
        if strings.HasSuffix(string(line), "@github.com") {
            u, _ := url.Parse(string(line))
            if password, ok := u.User.Password(); ok && strings.HasPrefix(password, "ghp_") {
                return password
            }
        }
    }

    return ""
}
