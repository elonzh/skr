package utils

import (
	"fmt"
	"net/url"
)

func DummyImage(width, height int, backgroundColor, foregroundColor, text, format string) string {
	return fmt.Sprintf("https://dummyimage.com/%dx%d/%s/%s.%s&text=%s", width, height, backgroundColor, foregroundColor, format, url.QueryEscape(text))
}
