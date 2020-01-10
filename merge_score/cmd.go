package merge_score

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand(_ *viper.Viper) *cobra.Command {
	var resultFilePath string
	cmd := &cobra.Command{
		Use:     "merge_score",
		Short:   "合并学生成绩单",
		Version: "v20200108",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return MergeScore(resultFilePath, args...)
		},
	}
	resultFilePathFlag := "resultFilePath"
	cmd.Flags().StringVarP(&resultFilePath, resultFilePathFlag, "p", "", "成绩汇总表路径, 成绩汇总表必须提供所有学生的姓名, 班级和科目信息")

	err := cmd.MarkFlagRequired(resultFilePathFlag)
	if err != nil {
		panic(err)
	}
	err = cmd.MarkFlagFilename(resultFilePathFlag, "xlsx")
	if err != nil {
		panic(err)
	}
	return cmd
}
