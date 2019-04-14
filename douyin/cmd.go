package douyin

import (
	"encoding/json"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var silent bool

var Cmd = &cobra.Command{
	Use:     "douyin",
	Short:   "爱抖音, 爱抖音小助手",
	Long:    `爱抖音小助手, 它能帮你解析抖音名片数据`,
	Version: "v20180716",
	Run: func(cmd *cobra.Command, args []string) {
		urls := viper.GetStringSlice("urls")
		logrus.WithField("urlsNum", len(urls)).Infoln("提取链接成功")

		users := make([]*User, 0, len(urls))
		for i, url := range urls {
			{
				logrus.Infof("开始请求第 %d 条链接: %s", i+1, url)
				user := User{URL: url}
				resp, _, errs := Fetch(user.URL)
				if len(errs) > 0 {
					logrus.Fatalln(errs)
				}
				doc, err := goquery.NewDocumentFromReader(resp.Body)
				if err != nil {
					logrus.WithError(err).Fatal("初始化页面失败")
				}
				Parse(doc, &user)
				users = append(users, &user)
			}
		}
		data, err := json.MarshalIndent(users, "", "  ")
		if err != nil {
			logrus.WithError(err).Fatal("数据编码失败")
		}
		fmt.Println(string(data))
	},
}

func init() {
	Cmd.PersistentFlags().BoolVar(&silent, "silent", false, "静默模式, 只在出错时输出日志")
	Cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "配置文件路径(默认为 config.yaml)")

	Cmd.PersistentFlags().StringSliceP("urls", "u", nil, "抖音分享链接")
	viper.BindPFlag("urls", Cmd.PersistentFlags().Lookup("urls"))
}
