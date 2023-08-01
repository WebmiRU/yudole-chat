package goodgame

type MessageIncomeData struct {
	ChannelId        string `json:"channel_id"`
	UserId           int    `json:"user_id"`
	UserName         string `json:"user_name"`
	MessageId        string `json:"message_id"`
	Timestamp        int    `json:"timestamp"`
	Text             string `json:"text"`
	ClientsInChannel int    `json:"clients_in_channel"`
	UsersInChannel   int    `json:"users_in_channel"`
	ErrorNum         int    `json:"error_num"`
	ErrorMsg         string `json:"errorMsg"`
}

type MessageIncome struct {
	Type string            `json:"type"`
	Data MessageIncomeData `json:"data"`
}

type MessageAuth struct {
	Type string          `json:"type"`
	Data MessageAuthData `json:"data"`
}

type MessageAuthData struct {
	SiteId int    `json:"site_id"`
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

type MessageJoin struct {
	Type string          `json:"type"`
	Data MessageJoinData `json:"data"`
}

type MessageJoinData struct {
	ChannelId string `json:"channel_id"`
	Hidden    bool   `json:"hidden"`
	Mobile    bool   `json:"mobile"`
}
