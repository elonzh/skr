package dingtalk

type Response struct {
	ErrorMessage string `json:"errmsg"`
	ErrorCode    int    `json:"errcode"`
}

type Message struct {
	Type       string      `json:"msgtype"`
	Text       *Text       `json:"text,omitempty"`
	Link       *Link       `json:"link,omitempty"`
	Markdown   *Markdown   `json:"markdown,omitempty"`
	ActionCard *ActionCard `json:"actionCard,omitempty"`
	FeedCard   *FeedCard   `json:"feedCard,omitempty"`
	At         *At         `json:"at,omitempty"`
}

type Text struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Link struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	MessageURL string `json:"messageUrl"`
	PictureURL string `json:"picUrl"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type At struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

type ActionCard struct {
	Title             string   `json:"title"`
	Text              string   `json:"text"`
	HideAvatar        string   `json:"hideAvatar"`
	ButtonOrientation string   `json:"btnOrientation"`
	Buttons           []Button `json:"btns,omitempty"`
	SingleTitle       string   `json:"singleTitle,omitempty"`
	SingleURL         string   `json:"singleURL"`
}

type Button struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}

type FeedCard struct {
	Links []FeedCardLink `json:"links"`
}

type FeedCardLink struct {
	Title      string `json:"title"`
	MessageURL string `json:"messageURL"`
	PictureURL string `json:"picURL"`
}
