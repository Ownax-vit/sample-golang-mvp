package services

import (
	"chat-project/internal/domain"
	"chat-project/internal/dto"
	"chat-project/internal/storage"
	"context"
	"fmt"
	"time"
)

type ChatService struct {
	chatRepo storage.ChatRepo
}

func New(chatRepo storage.ChatRepo) *ChatService {
	return &ChatService{
		chatRepo: chatRepo,
	}
}

// Создать чат
func (c ChatService) Create(ctx context.Context, chatIn dto.ChatIn) (*dto.ChatResponse, error) {
	chat := domain.Chat{
		Title:     chatIn.Title,
		CreatedAt: time.Now(),
	}

	chat, err := c.chatRepo.CreateChat(ctx, chat)
	if err != nil {
		return nil, fmt.Errorf("error while creating chat: %w", err)
	}

	return &dto.ChatResponse{
		Title:     chat.Title,
		ID:        chat.ID,
		CreatedAt: chat.CreatedAt.Format(time.RFC3339),
	}, nil
}

// Добавить сообщение в чат
func (c ChatService) AddMessage(ctx context.Context, chatId int, message dto.MessageIn) (*dto.MessageResponse, error) {
	msg, err := c.chatRepo.AddMessage(
		ctx,
		domain.Message{
			ChatId:    chatId,
			Text:      message.Text,
			CreatedAt: time.Now(),
		},
		chatId,
	)

	if err != nil && err == storage.ChatNotFoundError {
		return nil, fmt.Errorf("error while adding msg to chat with id %d: %w", chatId, err)
	} else if err != nil {
		return nil, fmt.Errorf("error while adding message: %w", err)
	}

	return &dto.MessageResponse{
		ID:        msg.ID,
		ChatId:    msg.ChatId,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt.Format(time.RFC3339),
	}, nil
}

// Получить чат с лимитом последних сообщений
func (c ChatService) GetWithMessages(ctx context.Context, chatId int) (*dto.ChatWithMessagesResponse, error) {
	chat, err := c.chatRepo.GetWithMessages(ctx, chatId)
	if err != nil && err == storage.ChatNotFoundError {
		return nil, fmt.Errorf("error while receive chat with id: %d : %w", chatId, err)
	} else if err != nil {
		return nil, fmt.Errorf("error while getting chat with messages: %w", err)
	}

	messages := make([]dto.MessageResponse, 0, len(chat.Messages))
	for _, msg := range chat.Messages {
		messages = append(messages, dto.MessageResponse{
			ID:        msg.ID,
			ChatId:    msg.ChatId,
			Text:      msg.Text,
			CreatedAt:msg.CreatedAt.Format(time.RFC3339),
		})
	}

	return &dto.ChatWithMessagesResponse{
		ID:        chat.ID,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt.Format(time.RFC3339),
		Messages:  messages,
	}, nil
}

// Удалить чат
func (c ChatService) DeleteChat(ctx context.Context, chatId int) error {
	err := c.chatRepo.DeleteChat(ctx, chatId)
	if err != nil && err == storage.ChatNotFoundError {
		return fmt.Errorf("chat with id %d not found: %w", chatId, err)
	} else if err != nil {
		return fmt.Errorf("error while deleting chat: %w", err)
	}

	return nil
}
