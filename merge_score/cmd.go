package merge_score

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand(_ *viper.Viper) *cobra.Command {
	var resultFilePath string
	cmd := &cobra.Command{
		Use:     "merge_score",
		Version: "v20200108",
		Short:   "合并学生成绩单",
		Long: `
			合并学生成绩至成绩汇总表

			数据表样例如下:
			  ___________________________________________________
			1|    国际学院20xx级xxx班20xx-20xx学年第x学期成绩统计   |
			2|    学生信息      |  语言课程分数   |  专业课程分数   |
			3| 序号 姓名 语言班 | 语言课1 语言课2 | 专业课1 专业课2 |
				A   B     C       D      E         F      G

			成绩汇总表需要提供所有学生姓名班级及科目信息，注意名字不要出现不一致的`,
		Example: `
			# 成绩表路径 可以为文件名或文件夹， 可以指定多个， 用空格隔开
			skr merge_score -p ".\2020第一学期成绩\2020传媒第一学期分数.xlsx" ".\2020第一学期成绩\专业\"  ".\2020第一学期成绩\语言\"`,
		Args: cobra.MinimumNArgs(1),
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
	return cmd
}
