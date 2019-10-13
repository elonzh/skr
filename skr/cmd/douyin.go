package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/earlzo/skr/douyin"
)

func newDouyinCommand(v *viper.Viper) *cobra.Command {
	var urls []string
	cmd := &cobra.Command{
		Use:     "douyin",
		Short:   "解析抖音名片数据",
		Version: "v20180716",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return douyin.Run(urls)
		},
	}
	cmd.Flags().StringSliceVarP(&urls, "urls", "u", nil, "抖音分享链接")
	var err error
	if err = v.BindPFlag(cmd.Name()+".urls", cmd.Flags().Lookup("urls")); err != nil {
		logrus.WithError(err).Fatalln()
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newDouyinCommand(v))
}
