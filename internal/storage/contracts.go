package storage

import (
	"chat-project/internal/domain"
	"context"
	"errors"
)

var ChatNotFoundError = errors.New("chat not found")

type ChatRepo interface {
	CreateChat(ctx context.Context, chat domain.Chat) (domain.Chat, error)
	GetChatByID(ctx context.Context, chatId int) (domain.Chat, error)
	AddMessage(ctx context.Context, msg domain.Message, chatId int) (domain.Message, error)
	GetWithMessages(ctx context.Context, chatId int) (domain.Chat, error)
	DeleteChat(ctx context.Context, chatId int) error
}

type ChatListener interface {
	Subscribe(ctx context.Context, chatId int) <-chan domain.Message
	Publish(ctx context.Context, chatId int, message domain.Message) error
}
