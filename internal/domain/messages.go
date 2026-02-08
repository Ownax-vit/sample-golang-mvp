package domain

import "time"

type Message struct {
	ID        int       `json:"Id"        example:"125216"`
	ChatId    int       `json:"ChatId"    example:"125216"`
	Text      string    `json:"Text"      example:"Hello world!"`
	CreatedAt time.Time `json:"createdAt"`
}
