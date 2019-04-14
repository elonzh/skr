package gaoxiaojob

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var storageFilename string
var keywords []string
var debug bool

var Cmd = &cobra.Command{
	Use:     "gaoxiaojob",
	Version: "v20190409",
	Short:   "抓取 高校人才网(http://gaoxiaojob.com/) 的最近招聘信息并根据关键词推送至钉钉",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := Run(storageFilename, args[0], keywords, debug); err != nil {
			logrus.WithError(err).Fatalln()
		}
	},
}

func init() {
	Cmd.Flags().StringArrayVarP(&keywords, "keywords", "k", []string{}, "关键词")
	Cmd.Flags().StringVarP(&storageFilename, "storage", "s", "storage.boltdb", "历史记录数据路径")
	Cmd.Flags().BoolVarP(&debug, "verbose", "v", false, "调试模式")
}
