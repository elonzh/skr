package cmd

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/earlzo/skr/douyin"
	"github.com/earlzo/skr/gaoxiaojob"
)

var cfgFile = ""
var rootCmd = &cobra.Command{
	Use:   "skr",
	Short: "🏁  skr~ skr~",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatalln()
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfgFile, "配置文件路径")

	var (
		err  error
		name string
	)
	name = "log-level"
	rootCmd.PersistentFlags().Uint32(name, uint32(logrus.InfoLevel), "")
	if err = viper.BindPFlag("logLevel", rootCmd.PersistentFlags().Lookup(name)); err != nil {
		logrus.WithError(err).Fatalln()
	}
	v := viper.GetViper()
	rootCmd.AddCommand(newDouyinCommand(v))
	rootCmd.AddCommand(newGaoxiaoJobCommand(v))
}

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

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		exe, err := os.Executable()
		if err != nil {
			logrus.WithError(err).Fatalln("获取程序路径失败")
		}
		cfgDir, err := filepath.Abs(filepath.Dir(exe))
		if err != nil {
			logrus.WithError(err).Fatalln("获取程序所在文件夹路径失败")
		}
		viper.AddConfigPath(cfgDir)
		viper.SetConfigName("skr")
	}
	if err := viper.ReadInConfig(); err == nil {
		logrus.WithField("cfgFile", viper.ConfigFileUsed()).Debugln("成功找到配置文件")
	} else {
		logrus.WithError(err).Debugln("没有找到配置文件")
	}
}

func initLogger() {
	level := logrus.Level(viper.GetInt("logLevel"))
	logrus.SetLevel(level)
	if level >= logrus.DebugLevel {
		rootCmd.DebugFlags()
		viper.Debug()
	}
}
