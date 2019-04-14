package douyin

import (
	"github.com/parnurzeal/gorequest"
)

// Fetch 拉取抖音用户名片页面
func Fetch(url string) (gorequest.Response, string, []error) {
	request := gorequest.New()
	return request.Get(url).Set(
		"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
	).End()
}
