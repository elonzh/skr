package cmd

import (
	"github.com/elonzh/skr/winfocus"
)

func init() {
	rootCmd.AddCommand(winfocus.NewCommand(v))
}
