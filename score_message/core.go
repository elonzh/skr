package score_message

import (
	"fmt"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/sirupsen/logrus"

	"github.com/elonzh/skr/merge_score"
	"github.com/elonzh/skr/pkg/utils"
)

func RenderMessage(file *excelize.File, skipRows uint32) error {
	rows, err := file.GetRows(file.GetSheetName(1))
	if err != nil {
		return err
	}
	newFile := excelize.NewFile()
	sheetName := "Sheet1"
	rows = rows[skipRows:]
	const nameField = "姓名"
	headers, maps := merge_score.ReadMap(rows, nil, []string{"序号", "语言班"})
	logrus.WithFields(logrus.Fields{
		"rows":    rows,
		"headers": headers,
		"maps":    maps,
	}).Debugln("数据加载完成")
	for i, d := range maps {
		name := d[nameField]
		s := strings.Builder{}
		s.WriteString(fmt.Sprintf("%s家长你好，以下是%s本学期的分数：\n", name, name))
		for idx, h := range headers {
			if h != nameField && d[h] != "" {
				s.WriteString(fmt.Sprintf("%s %s分", h, d[h]))
			}
			if idx != len(headers)-1 {
				s.WriteString("；")
			} else {
				s.WriteString("。")
			}
		}
		s.WriteString("以上课程满分均为100分。")
		err := newFile.SetCellStr(sheetName, merge_score.PointToAxis(i+1, 1), s.String())
		if err != nil {
			return err
		}
	}
	outputPath := utils.PrefixedPath("信息生成表-", file.Path)
	logrus.WithField("信息生成表位置", outputPath).Info("信息生成表生成完成")
	return newFile.SaveAs(outputPath)
}
