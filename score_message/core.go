package score_message

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

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
		needRetake := false
		name := d[nameField]
		buf.WriteString(fmt.Sprintf("%s家长你好，以下是%s本学期的分数：\n", name, name))
		for idx, h := range headers {
			if h != nameField && d[h] != "" {
				if IsDigit(d[h]) {
					buf.WriteString(fmt.Sprintf("%s %s分", h, d[h]))
				} else {
					buf.WriteString(fmt.Sprintf("%s %s", h, d[h]))
					if strings.Contains(d[h], "重修") {
						needRetake = true
					}
				}
				if idx != len(headers)-1 {
					buf.WriteString("；")
				} else {
					buf.WriteString("。")
				}
			}
		}
		buf.WriteString("以上课程满分均为100分。")
		if needRetake {
			buf.WriteString(fmt.Sprintf("“重修”课程是由于%s在本学期的到课率少于总课程的三分之二，根据我院学生手册规定，需要进行重修。具体重修事宜请学生本人等待辅导员通知。", name))
		}
		buf.WriteString("请家长提醒孩子，在班级群中查看今年的暑假作业要求，按时完成作业。\n\n")
	}

	outputPath := filepath.Join(filepath.Dir(file.Path), "成绩信息.txt")
	err = ioutil.WriteFile(outputPath, buf.Bytes(), 0666)
	if err != nil {
		return err
	}
	logrus.WithField("成绩信息位置", outputPath).Info("成绩信息生成完成")
	return nil
}
