package cmd

import (
	"github.com/earlzo/skr/douyin"
	"github.com/earlzo/skr/gaoxiaojob"
)

func init() {
	rootCmd.AddCommand(douyin.Cmd)
	rootCmd.AddCommand(gaoxiaojob.Cmd)
}
