package douyin

import (
	"bytes"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/gobuffalo/packr"
	"github.com/golang/freetype/truetype"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

var font *truetype.Font

func init() {
	var err error
	box := packr.NewBox("./static")
	data, err := box.Find("iconfont_da2e2ef.ttf")
	if err != nil {
		logrus.WithError(err).Panicln("载入字体数据失败")
	}
	font, err = truetype.Parse(data)
	if err != nil {
		logrus.WithError(err).Panicln("载入字体数据失败")
	}
	logrus.WithFields(logrus.Fields{
		"NameIDFontFamily":       font.Name(truetype.NameIDFontFamily),
		"NameIDFontFullName":     font.Name(truetype.NameIDFontFullName),
		"NameIDNameTableVersion": font.Name(truetype.NameIDNameTableVersion),
		"NameIDCopyright":        font.Name(truetype.NameIDCopyright),
	}).Debug("载入字体数据成功")
}

// 一个数字可能对应多个编码, 但相应地 Index 是一致的
var numMap = map[uint16]int{
	3:  0,
	2:  1,
	5:  2,
	4:  3,
	6:  4,
	7:  5,
	8:  6,
	10: 7,
	11: 8,
	9:  9,
}

var runeMap = map[int]rune{
	0: '0',
	1: '1',
	2: '2',
	3: '3',
	4: '4',
	5: '5',
	6: '6',
	7: '7',
	8: '8',
	9: '9',
}

func cleanStr(str string) string {
	return strings.Replace(str, " ", "", -1)
}

func parseNumStr(numBlock *goquery.Selection) string {
	buf := bytes.NewBuffer(nil)
	for _, node := range numBlock.Contents().Nodes {
		if node.Type == html.TextNode {
			buf.WriteString(cleanStr(node.Data))
		} else {
			doc := goquery.NewDocumentFromNode(node)
			num := cleanStr(doc.Text())
			for _, r := range num {
				n := numMap[uint16(font.Index(r))]
				r = runeMap[n]
				buf.WriteRune(r)
			}
		}
	}
	numStr := buf.String()
	return numStr
}

func numStr2num(numStr string) uint {
	var err error
	var num uint
	if strings.HasSuffix(numStr, "w") {
		var numFloat float64
		numFloat, err = strconv.ParseFloat(numStr[:len(numStr)-utf8.RuneLen('w')], 32)
		num = uint(math.Round(numFloat * 10000))
	} else {
		var numInt int
		numInt, err = strconv.Atoi(numStr)
		num = uint(numInt)
	}
	if err != nil {
		logrus.WithField("numStr", numStr).WithError(err).Panicln("解析数字失败")
	}
	return num
}

// Parse 解析抖音用户名片数据
func Parse(doc *goquery.Document, user *User) {
	// doc.Find(".shortid .iconfont").Each(func(i int, s *goquery.Selection) {
	//	t := strings.Replace(s.Text(), " ", "", -1)
	//	for _, r := range t {
	//		user.IdRunes = append(user.IdRunes, r)
	//		// fmt.Printf("%+q %d\n", r, user.IDs[i])
	//		// fmt.Println(id[i], idx, id[i] == rune_map[r], id[i] == num_map[idx])
	//	}
	// })

	userInfoDoc := doc.Find("#pagelet-user-info").First()
	personalCard := userInfoDoc.Find(".personal-card").First()

	avatar := personalCard.Find("span.author > img.avatar").First()
	val, _ := avatar.Attr("src")
	user.Avatar = val

	nickname := personalCard.Find("p.nickname").First()
	user.NickName = nickname.Text()

	idStr := personalCard.Find("p.shortid").First()
	user.ID = strings.TrimLeft(parseNumStr(idStr), "抖音ID：")

	signature := personalCard.Find("p.signature").First()
	user.Signature = signature.Text()

	focus := doc.Find(".follow-info > .focus > span.num").First()
	user.FocusNumStr = parseNumStr(focus)
	user.FocusNum = numStr2num(user.FocusNumStr)

	follower := doc.Find(".follower > span.num").First()
	user.FollowerNumStr = parseNumStr(follower)
	user.FollowerNum = numStr2num(user.FollowerNumStr)

	liked := doc.Find(".liked-num > span.num").First()
	user.LikesNumStr = parseNumStr(liked)
	user.LikesNum = numStr2num(user.LikesNumStr)

	extraInfo := personalCard.Find("p.extra-info").First()

	location := extraInfo.Find("span.location").First()
	user.Location = location.Text()

	constellation := extraInfo.Find("span.constellation").First()
	user.Constellation = constellation.Text()

	videoTab := userInfoDoc.Find(".video-tab").First()

	postNum := videoTab.Find(".user-tab .num")
	user.PostNumStr = parseNumStr(postNum)
	user.PostNum = numStr2num(user.PostNumStr)

	likedNum := videoTab.Find(".like-tab .num")
	user.LikedNumStr = parseNumStr(likedNum)
	user.LikedNum = numStr2num(user.LikedNumStr)
}
