package douyin

import "time"

// User 抖音用户信息
type User struct {
	URL string

	ID            string
	Avatar        string
	NickName      string
	Signature     string
	Location      string
	Constellation string

	FollowerNumStr string
	FollowerNum    uint
	FocusNumStr    string
	FocusNum       uint
	LikesNumStr    string
	LikesNum       uint

	PostNumStr  string
	PostNum     uint
	LikedNumStr string
	LikedNum    uint

	Time time.Time
}
