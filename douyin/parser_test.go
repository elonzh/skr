package douyin

import (
	"testing"
)

func TestParseNumStr(t *testing.T) {

}

func TestNumStr2Num(t *testing.T) {
	for _, c := range []struct {
		NumStr string
		Num    uint
	}{
		{"18.1w", 181000},
		{"0.1w", 1000},
		{"123", 123},
		{"1", 1},
	} {
		n := numStr2num(c.NumStr)
		if n != c.Num {
			t.Errorf("NumStr %s, Num %d, expected %d", c.NumStr, n, c.Num)
		}
	}
}
