package trovo

type MessageResponseUsers struct {
	Total int                   `json:"total"`
	Users []MessageResponseUser `json:"users"`
}

type MessageResponseUser struct {
	UserId    string `json:"user_id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	ChannelId string `json:"channel_id"`
}

type MessageResponseToken struct {
	Token string `json:"token"`
}

type MessageAuth struct {
	Type  string          `json:"type"`
	Nonce string          `json:"nonce"`
	Data  MessageAuthData `json:"data"`
}

type MessageAuthData struct {
	Token string `json:"token"`
}

type Message struct {
	Type        string             `json:"type"`
	Nonce       string             `json:"nonce"`
	ChannelInfo MessageChannelInfo `json:"channel_info"`
	Data        MessageData        `json:"data"`
}

type MessageChannelInfo struct {
	ChannelId string `json:"channel_id"`
}

type MessageData struct {
	Eid   string            `json:"eid"`
	Chats []MessageDataChat `json:"chats"`
	Gap   int               `json:"gap"` // for PONG message type
}

type MessageDataChat struct {
	Type     int      `json:"type"`
	Content  string   `json:"content"`
	NickName string   `json:"nick_name"`
	Avatar   string   `json:"avatar"`
	SubLv    string   `json:"sub_lv"`
	SubTier  string   `json:"sub_tier"`
	Medals   []string `json:"medals"`
	//Decos       interface{} `json:"decos"`
	Roles     []string `json:"roles"`
	MessageId string   `json:"message_id"`
	SenderId  int      `json:"sender_id"`
	SendTime  int      `json:"send_time"`
	//Uid         interface{} `json:"uid"`
	UserName string `json:"user_name"`
	//ContentData interface{} `json:"content_data"`
}

type MessagePing struct {
	Type  string `json:"type"`
	Nonce string `json:"nonce"`
}
