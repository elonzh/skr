package gaoxiaojob

import "time"

type Job struct {
	URL   string
	Title string

	Meta map[string]string
	// 分类
	Categories []string
	// 需求学科
	Subjects []string
	// 所属省份
	Provinces []string
	// 工作地点
	Locations []string

	// 发布时间
	PublishedAt *time.Time
	// 截止日期
	ExpireAt *time.Time

	Body string
}
