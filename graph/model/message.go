package model

type Message struct {
	ID        string `json:"id"`
	RoomID    string `json:"roomId"`
	Sender    *User  `json:"sender"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type OtherMessage struct {
	ID        int64  `json:"chat_id"` // Change this from string to int64
	RoomID    string `json:"room_id"`
	Sender    string `json:"sender_id"`
	Content   string `json:"text"`
	Timestamp string `json:"timestamp"`
}
