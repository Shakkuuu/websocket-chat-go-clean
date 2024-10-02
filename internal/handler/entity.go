package handler

// HTMLテンプレートに渡すためのデータ
type Data struct {
	Name    string
	Message string
}

// ユーザー名送信用
type SentUser struct {
	Name string `json:"name"`
}

// クライアントサーバ間でやりとりするメッセージ
type Message struct {
	RoomID      string   `json:"roomid"`
	Message     string   `json:"message"`
	Name        string   `json:"name"`
	ToName      string   `json:"toname"`
	AllUsers    []string `json:"allusers"`
	OnlineUsers []string `json:"onlineusers"`
}

// ルーム一覧送信用
type SentRoomsList struct {
	RoomsList []string `json:"roomslist"`
}
