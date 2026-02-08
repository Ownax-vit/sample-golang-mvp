package memory

import (
	"chat-project/internal/domain"
	"chat-project/internal/storage"
	"context"

	"github.com/google/uuid"
)

type ChatRepoMemory struct {
	chats map[int]domain.Chat
}

func NewUserRepoMemory() *ChatRepoMemory {
	return &ChatRepoMemory{
		chats: make(map[int]domain.Chat),
	}
}

func (r *ChatRepoMemory) CreateChat(ctx context.Context, chat domain.Chat) (domain.Chat, error) {
	chat.ID = int(uuid.New().ID())
	r.chats[chat.ID] = chat
	return chat, nil
}

func (r *ChatRepoMemory) GetChatByID(ctx context.Context, chatId int) (domain.Chat, error) {
	chat, exists := r.chats[chatId]
	if !exists {
		return domain.Chat{}, storage.ChatNotFoundError
	}
	return chat, nil
}

func (r *ChatRepoMemory) AddMessage(ctx context.Context, message domain.Message, chatId int) (domain.Message, error) {
	chat, exists := r.chats[chatId]
	if !exists {
		return domain.Message{}, storage.ChatNotFoundError
	}

	message.ID = int(uuid.New().ID())
	chat.Messages = append(chat.Messages, message)
	r.chats[chatId] = chat
	return message, nil
}

func (r *ChatRepoMemory) GetWithMessages(ctx context.Context, chatId int) (domain.Chat, error) {
	chat, exists := r.chats[chatId]
	if !exists {
		return domain.Chat{}, storage.ChatNotFoundError
	}
	return chat, nil
}

func (r *ChatRepoMemory) DeleteChat(ctx context.Context, chatId int) error {
	_, exists := r.chats[chatId]
	if !exists {
		return storage.ChatNotFoundError
	}
	delete(r.chats, chatId)
	return nil
}
