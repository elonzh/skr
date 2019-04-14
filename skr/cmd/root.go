package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
	//todo: config file support
	//cobra.OnInitialize(initConfig)
}

//func initConfig() {
//	if silent {
//		logrus.SetLevel(logrus.ErrorLevel)
//		logrus.SetLevel(logrus.ErrorLevel)
//	}
//	if cfgFile != "" {
//		viper.SetConfigFile(cfgFile)
//	} else {
//		exe, err := os.Executable()
//		if err != nil {
//			logrus.WithError(err).Fatalln("获取程序路径失败")
//		}
//		cfgDir, err := filepath.Abs(filepath.Dir(exe))
//		if err != nil {
//			logrus.WithError(err).Fatalln("获取程序所在文件夹路径失败")
//		}
//		viper.AddConfigPath(cfgDir)
//		viper.SetConfigName("config")
//	}
//
//	viper.AutomaticEnv()
//	if err := viper.ReadInConfig(); err == nil {
//		logrus.WithField("cfgFile", viper.ConfigFileUsed()).Infoln("成功找到配置文件")
//	} else {
//		logrus.WithError(err).Fatalln("读取配置文件失败")
//	}
//}
