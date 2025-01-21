package model

type Message struct {
	ID        string `json:"id"`
	RoomID    string `json:"roomId"`
	Sender    *User  `json:"sender"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}
