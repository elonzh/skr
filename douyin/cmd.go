package douyin

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func NewCommand(v *viper.Viper) *cobra.Command {
	var urls []string
	cmd := &cobra.Command{
		Use:     "douyin",
		Short:   "解析抖音名片数据",
		Version: "v20180716",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(urls)
		},
	}
	cmd.Flags().StringSliceVarP(&urls, "urls", "u", nil, "抖音分享链接")
	var err error
	if err = v.BindPFlag(cmd.Name()+".urls", cmd.Flags().Lookup("urls")); err != nil {
		logrus.WithError(err).Fatalln()
	}
	return cmd
}
