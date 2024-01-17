package common

import (
    "errors"
    "fmt"
    "github.com/avast/retry-go"
    "github.com/sirupsen/logrus"
    "io"
    "net/http"
    "time"
)

func Get(url string) ([]byte, error) {
    var b []byte

    err := retry.Do(
        func() error {
            eclint := http.Client{
                Timeout: 5 * time.Second,
            }

            response, err1 := eclint.Get(url)
            if err1 != nil {
                return errors.New("请求错误：" + err1.Error())
            }

            defer func() {
                _ = response.Body.Close()
            }()

            b, err1 = io.ReadAll(response.Body)
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
