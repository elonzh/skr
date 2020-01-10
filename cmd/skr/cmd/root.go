package cmd

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/elonzh/skr/pkg/utils"
)

var (
	cfgFile = ""
	rootCmd = &cobra.Command{
		Use:   "skr",
		Short: "🏁 skr~ skr~",
	}
	v = viper.GetViper()
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatalln()
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogger, func() {
		utils.NormalizeAll(rootCmd)
	})
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
