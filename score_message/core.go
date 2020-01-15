package score_message

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/sirupsen/logrus"

	"github.com/elonzh/skr/merge_score"
)

var digitPattern = regexp.MustCompile(`^[\d\.]+$`)

func IsDigit(s string) bool {
	return digitPattern.MatchString(s)
}

func RenderMessage(file *excelize.File, skipRows uint32) error {
	rows, err := file.GetRows(file.GetSheetName(1))
	if err != nil {
		return err
	}
	rows = rows[skipRows:]
	const nameField = "姓名"
	headers, maps := merge_score.ReadMap(rows, nil, []string{"序号", "语言班"})
	logrus.WithFields(logrus.Fields{
		"rows":    rows,
		"headers": headers,
		"maps":    maps,
	}).Debugln("数据加载完成")
	buf := bytes.Buffer{}
	for _, d := range maps {
		name := d[nameField]
		buf.WriteString(fmt.Sprintf("%s家长你好，以下是%s本学期的分数：\n", name, name))
		for idx, h := range headers {
			if h != nameField && d[h] != "" {
				if IsDigit(d[h]) {
					buf.WriteString(fmt.Sprintf("%s %s分", h, d[h]))
				} else {
					buf.WriteString(fmt.Sprintf("%s %s", h, d[h]))
				}
				if idx != len(headers)-1 {
					buf.WriteString("；")
				} else {
					buf.WriteString("。")
				}
			}
		}
		buf.WriteString("以上课程满分均为100分。\n\n")
	}

	outputPath := filepath.Join(filepath.Dir(file.Path), "成绩信息.txt")
	err = ioutil.WriteFile(outputPath, buf.Bytes(), 0666)
	if err != nil {
		return err
	}
	logrus.WithField("成绩信息位置", outputPath).Info("成绩信息生成完成")
	return nil
}
