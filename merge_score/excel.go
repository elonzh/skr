package merge_score

import (
	"fmt"
	"sort"
	"strings"

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

func ReadMap(rows [][]string, headers, excludeHeaders []string) ([]string, []map[string]string) {
	headerMap := make(map[string]int)
	for i, v := range rows[0] {
		h := strings.TrimSpace(v)
		if h != "" {
			headerMap[h] = i
		}
	}
	if len(headers) > 0 {
		keep := make(map[string]struct{}, len(headers))
		for _, h := range headers {
			keep[h] = struct{}{}
		}
		for k := range headerMap {
			if _, ok := keep[k]; !ok {
				delete(headerMap, k)
			}
		}
	}
	if len(excludeHeaders) > 0 {
		for _, k := range excludeHeaders {
			delete(headerMap, k)
		}
	}
	result := make([]map[string]string, 0, len(rows[1:]))
	for _, row := range rows[1:] {
		m := make(map[string]string, len(headerMap))
		for key, index := range headerMap {
			if index >= len(row) {
				continue
			}
			m[key] = row[index]
		}
		result = append(result, m)
	}
	cleanHeaders := make([]string, len(headerMap))
	for k := range headerMap {
		cleanHeaders = append(cleanHeaders, k)
	}
	sort.Slice(cleanHeaders, func(i, j int) bool {
		return headerMap[cleanHeaders[i]] < headerMap[cleanHeaders[j]]
	})
	return cleanHeaders, result
}
