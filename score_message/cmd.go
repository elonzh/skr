package score_message

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand(_ *viper.Viper) *cobra.Command {
	var skipRows uint32
	cmd := &cobra.Command{
		Use:     "score_message",
		Version: "v20200111",
		Short:   "根据学生成绩生成信息",
		Example: `skr score_message ".\2020第一学期成绩\已汇总-2020传媒第一学期分数.xlsx"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := excelize.OpenFile(args[0])
			if err != nil {
				return err
			}
			return RenderMessage(file, skipRows)
		},
	}
	cmd.Flags().Uint32Var(&skipRows, "skipRows", 2, "跳过 N 行，从 N + 1 行处理")
	return cmd
}
