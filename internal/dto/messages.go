package dto

type MessageIn struct {
	Text    string `json:"Text"      example:"Hello world!"`
}

type MessageResponse struct {
	ID        int    `json:"Id"        example:"125216"`
	ChatId   int    `json:"ChatId"    example:"125216"`
	Text      string `json:"Text"      example:"Hello world!"`
	CreatedAt string `json:"createdAt" example:"2024-01-01T12:00:00Z"`
}

type MessagesResponse struct {
	Messages []MessageResponse `json:"messages"`
}
