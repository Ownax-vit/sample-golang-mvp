package dto

type ChatIn struct {
	Title string `json:"Title"    example:"Тестовый чат"`
}

type ChatResponse struct {
	Title     string `json:"Title"    example:"Тестовый чат"`
	ID        int    `json:"Id"       example:"125216"`
	CreatedAt string `json:"CreatedAt" example:"2024-01-01T12:00:00Z"`
}

type ChatWithMessagesResponse struct {
	Title     string            `json:"Title" example:"Тестовый чат"`
	ID        int               `json:"Id" example:"125216"`
	CreatedAt string            `json:"CreatedAt" example:"2024-01-01T12:00:00Z"`
	Messages  []MessageResponse `json:"messages"`
}

type ChatsResponse struct {
	Chats []ChatResponse `json:"chats"`
}
