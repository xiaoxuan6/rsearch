package command

import (
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
	files, err := fetchFiles()
	if err != nil {
		return err
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
	dir, err := os.Getwd()
	if err != nil {
		return files, err
	}

	path := dir + common.RepositoryPath
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
