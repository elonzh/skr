package cmd

import (
	"github.com/elonzh/skr/douyin"
	"github.com/elonzh/skr/gaoxiaojob"
	"github.com/elonzh/skr/merge_score"
	"github.com/elonzh/skr/score_message"
	"github.com/elonzh/skr/slot_machine"
)

func init() {
	rootCmd.AddCommand(douyin.NewCommand(v))
	rootCmd.AddCommand(gaoxiaojob.NewCommand(v))
	rootCmd.AddCommand(merge_score.NewCommand(v))
	rootCmd.AddCommand(score_message.NewCommand(v))
	rootCmd.AddCommand(slot_machine.NewCommand(v))
}
