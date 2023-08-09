package messages

type System struct {
	Type    string `json:"type"`
	Service string `json:"service"`
	User    User   `json:"user"`
	Channel string `json:"channel"`
	Value   string `json:"value"`
}
