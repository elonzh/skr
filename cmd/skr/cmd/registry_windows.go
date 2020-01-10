package cmd

func init() {
	rootCmd.AddCommand(winfocus.NewCommand(v))
}
