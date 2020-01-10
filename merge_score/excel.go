package merge_score

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func PointToAxis(x, y int) string {
	return fmt.Sprintf("%s%d", MustColumnNumberToName(y), x)
}
func MustColumnNumberToName(num int) string {
	name, err := excelize.ColumnNumberToName(num)
	if err != nil {
		panic(err)
	}
	return name
}
