package merge_score

import (
	"errors"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/sirupsen/logrus"
)

//  数据表样例如下:
//   ___________________________________________________
// 1|    国际学院20xx级xxx班20xx-20xx学年第x学期成绩统计   |
// 2|    学生信息      |  语言课程分数   |  专业课程分数   |
// 3| 序号 姓名 语言班 | 语言课1 语言课2 | 专业课1 专业课2 |
//     A   B     C       D      E         F      G
func LoadScoreTable(file *excelize.File, skipRows int) (*ScoreTable, error) {
	logrus.WithFields(logrus.Fields{
		"SheetMap": file.GetSheetMap(),
		"Path":     file.Path,
	}).Debugln("开始加载数据")
	rows, err := file.GetRows(file.GetSheetMap()[1])
	if err != nil {
		return nil, err
	}
	rows = rows[skipRows:]
	headers := rows[0]
	for i, v := range headers {
		headers[i] = strings.TrimSpace(v)
	}
	if len(headers) <= 4 {
		logrus.Warnln("列数小于 4 列，是不是没有数据？")
	}

	t := &ScoreTable{
		file:     file,
		SkipRows: skipRows,
		ScoreMap: make(map[StudentSubject]StudentSubjectScore, len(rows)-1),
	}
	const subjectNameIndexOffset = 3
	for rowIndex, row := range rows[1:] {
		x := rowIndex + skipRows + 1 + 1
		className := strings.TrimSpace(row[2])
		studentName := strings.TrimSpace(row[1])
		if className == "" || studentName == "" {
			logrus.WithFields(logrus.Fields{
				"行":  x,
				"班级": className,
				"学生": studentName,
			}).Warnln("学生信息缺失, 跳过该行")
			continue
		}
		for subjectNameIndex, rawScoreStr := range row[3:] {
			y := subjectNameIndex + subjectNameIndexOffset + 1
			subjectName := headers[y-1]
			if subjectName == "" {
				logrus.WithFields(logrus.Fields{
					"行":  x,
					"列":  MustColumnNumberToName(y),
					"班级": className,
					"学生": studentName,
				}).Warnln("科目缺失, 跳过该列")
				continue
			}
			rawScoreStr = strings.TrimSpace(rawScoreStr)
			if rawScoreStr == "" {
				t.ScoreMap[StudentSubject{
					ClassName:   className,
					StudentName: studentName,
					SubjectName: headers[subjectNameIndex+subjectNameIndexOffset],
				}] = StudentSubjectScore{
					Score: -1,
					X:     x,
					Y:     y,
				}
				logrus.WithFields(logrus.Fields{
					"行":    x,
					"列":    MustColumnNumberToName(y),
					"班级":   className,
					"学生":   studentName,
					"科目":   subjectName,
					"分数数据": rawScoreStr,
				}).Warnln("分数数据为空")
			} else {
				rawScore, err := strconv.ParseFloat(rawScoreStr, 64)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"行":    x,
						"列":    MustColumnNumberToName(y),
						"班级":   className,
						"学生":   studentName,
						"科目":   subjectName,
						"分数数据": rawScoreStr,
					}).Warnln("错误的分数数据, 跳过该列")
				} else {
					t.ScoreMap[StudentSubject{
						ClassName:   className,
						StudentName: studentName,
						SubjectName: subjectName,
					}] = StudentSubjectScore{
						Score: rawScore,
						X:     x,
						Y:     y,
					}
					logrus.WithFields(logrus.Fields{
						"行":    x,
						"列":    MustColumnNumberToName(y),
						"班级":   className,
						"学生":   studentName,
						"科目":   subjectName,
						"分数数据": rawScoreStr,
					}).Infoln("分数加载成功")
				}
			}

		}
	}
	return t, nil
}

func MergeScore(resultFilePath string, scoreFilePaths ...string) error {
	if len(scoreFilePaths) == 0 {
		return errors.New("should provide at least one scoreFilePath")
	}
	const skipRows = 2

	file, err := excelize.OpenFile(resultFilePath)
	if err != nil {
		return err
	}
	t, err := LoadScoreTable(file, skipRows)
	if err != nil {
		return err
	}
	sheetName := file.GetSheetName(1)

	scoreTables := make([]*ScoreTable, 0, len(scoreFilePaths))
	for _, scoreFilePath := range scoreFilePaths {
		file, err := excelize.OpenFile(scoreFilePath)
		if err != nil {
			return err
		}
		t, err := LoadScoreTable(file, skipRows)
		if err != nil {
			return err
		}
		scoreTables = append(scoreTables, t)
	}

	for subject, score := range t.ScoreMap {
		for _, scoreTable := range scoreTables {
			realScore, ok := scoreTable.ScoreMap[subject]
			if ok && realScore.Score != -1 {
				s := StudentSubjectScore{
					Score: realScore.Score,
					X:     score.X,
					Y:     score.Y,
				}
				t.ScoreMap[subject] = s
				if err := file.SetCellFloat(sheetName, s.GetAxis(), s.Score, 2, 64); err != nil {
					logrus.WithError(err).Warningln("分数设置错误")
				}
				break
			}
		}
	}
	outputPath := "已汇总-" + file.Path
	err = file.SaveAs(outputPath)
	if err != nil {
		return err
	}
	logrus.WithField("汇总表位置", outputPath).Info("成绩汇总完成")
	return nil
}