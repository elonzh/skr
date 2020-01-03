package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/elonzh/skr/gaoxiaojob"
)

func newGaoxiaoJobCommand(v *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gaoxiaojob",
		Version: "v20190409",
		Short:   "抓取 高校人才网(http://gaoxiaojob.com/) 的最近招聘信息并根据关键词推送至钉钉",
		Args:    cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			v.Set(cmd.Name()+".webhookURL", args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return gaoxiaojob.Run(v.GetString(cmd.Name()+".storage"), args[0], v.GetStringSlice(cmd.Name()+".keywords"))
		},
	}
	var err error
	cmd.Flags().StringArrayP("keywords", "k", nil, "关键词")
	if err = v.BindPFlag(cmd.Name()+".keywords", cmd.Flags().Lookup("keywords")); err != nil {
		return nil
	}
	cmd.Flags().StringP("storage", "s", "storage.boltdb", "历史记录数据路径")
	if err = v.BindPFlag(cmd.Name()+".storage", cmd.Flags().Lookup("storage")); err != nil {
		return nil
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newGaoxiaoJobCommand(v))
}
