package winfocus

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/elonzh/skr/pkg/utils"
)

func NewCommand(v *viper.Viper) *cobra.Command {
	currentDir, err := filepath.Abs("")
	if err != nil {
		logrus.WithError(err).Fatalln()
	}
	outputDir := filepath.Join(currentDir, "Windows Focus")
	cmd := &cobra.Command{
		Use:     "winfocus",
		Version: "v20191012",
		Short:   "导出 Windows 聚焦图片",
		RunE: func(cmd *cobra.Command, args []string) error {
			localAppDataPath := os.Getenv("localappdata")
			assetPath := filepath.Join(localAppDataPath, "/Packages/Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy/LocalState/Assets")
			err := os.MkdirAll(outputDir, 0666)
			if err != nil {
				logrus.WithError(err).WithField("dir", outputDir).Fatalln("make dir failed")
			} else {
				logrus.WithField("outputDir", outputDir).Infoln("outputDir created")
			}
			return filepath.Walk(assetPath, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return err
				}
				filename := info.Name() + ".jpg"
				dst := filepath.Join(outputDir, filename)
				err = utils.CopyFile(dst, path, 0666)
				if err != nil {
					logrus.WithError(err).WithFields(logrus.Fields{
						"src": path,
						"dst": dst,
					}).Warnln("fail to copy file")
				}
				return err
			})
		},
	}
	cmd.Flags().StringVarP(&outputDir, "output", "o", outputDir, "导出路径")
	return cmd
}
