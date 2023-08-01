package goodgame

//type Smiles struct {
//	Smiles []Smile `json:"smiles"`
//}

type Smile struct {
	Id         string `json:"id"`
	Key        string `json:"key"`
	Level      int    `json:"level"`
	Paid       string `json:"paid"`
	Bind       string `json:"bind"`
	InternalId int    `json:"internal_id"`
	ChannelId  int    `json:"channel_id"`
	Channel    string `json:"channel"`
	Nickname   string `json:"nickname"`
	Donat      int    `json:"donat"`
	Premium    int    `json:"premium"`
	Animated   int    `json:"animated"`
	Images     Images `json:"images"`
}

type Images struct {
	Small string `json:"small"`
	Big   string `json:"big"`
	Gif   string `json:"gif"`
}
