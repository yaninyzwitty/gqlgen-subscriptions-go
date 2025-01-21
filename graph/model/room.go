package model

type Room struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Participants []*User    `json:"participants"`
	CreatedAt    string     `json:"createdAt"`
	Messages     []*Message `json:"messages"`
}
