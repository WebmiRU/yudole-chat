package messages

type User struct {
	Login     string `json:"login"`
	Nick      string `json:"nick"`
	AvatarUrl string `json:"avatar_url"`
	Color     string `json:"color"`
}

type Message struct {
	Text string `json:"text"`
	Html string `json:"html"`
	//SmileLess string `json:"smile_less"`
}

type Channel struct {
	Service string  `json:"service"`
	Type    string  `json:"type"`
	User    User    `json:"user"`
	Message Message `json:"message"`
}
