package merge_score

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type StudentSubject struct {
	ClassName   string
	StudentName string
	SubjectName string
}

func (s *StudentSubject) String() string {
	return fmt.Sprintf("%s %s %s", s.ClassName, s.StudentName, s.SubjectName)
}

type StudentSubjectScore struct {
	ScoreData string
	X, Y      int
}

func (s *StudentSubjectScore) GetAxis() string {
	return PointToAxis(s.X, s.Y)
}

func (s *StudentSubjectScore) String() string {
	return fmt.Sprintf("score %s at %s", s.ScoreData, s.GetAxis())
}

type ScoreTable struct {
	file     *excelize.File
	SkipRows uint32
	ScoreMap map[StudentSubject]StudentSubjectScore
}

func (t *ScoreTable) String() string {
	return fmt.Sprintf("%s\n%v", t.file.Path, t.ScoreMap)
}
