package cmd

import (
	"github.com/elonzh/skr/douyin"
	"github.com/elonzh/skr/gaoxiaojob"
	"github.com/elonzh/skr/merge_score"
)

func init() {
	rootCmd.AddCommand(douyin.NewCommand(v))
	rootCmd.AddCommand(gaoxiaojob.NewCommand(v))
	rootCmd.AddCommand(merge_score.NewCommand(v))
}
