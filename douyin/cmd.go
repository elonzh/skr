package douyin

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

func Run(urls []string) error {
	logrus.WithField("urlsNum", len(urls)).Infoln("提取链接成功")
	c := colly.NewCollector()
	users := make([]*User, 0, len(urls))
	c.OnRequest(func(request *colly.Request) {
		logrus.Debugf("开始请求链接: %s", request.URL)
		request.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	})
	c.OnResponse(func(response *colly.Response) {
		user := User{URL: response.Request.URL.String()}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response.Body))
		if err != nil {
			logrus.WithError(err).Fatal("初始化页面失败")
		}
		Parse(doc, &user)
		users = append(users, &user)
	})
	for _, url := range urls {
		if err := c.Visit(url); err != nil {
			logrus.WithError(err).Warningln("请求失败:", url)
		}
	}
	c.Wait()

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		logrus.WithError(err).Fatal("数据编码失败")
	}
	fmt.Println(string(data))
	return err
}
