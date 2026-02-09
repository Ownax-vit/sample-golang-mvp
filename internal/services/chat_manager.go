package services

import (
	"chat-project/internal/storage"
	"context"
)

// Менеджер слушателей, который будет хранить всех слушателей для разных чатов и создавать новых при необходимости
type ChatListenerManager struct {
	repoChat      storage.ChatRepo
	listener      storage.ChatListener
	chatListeners map[int]*ChatListener
}

func NewChatListenerManager(repoChat storage.ChatRepo, listener storage.ChatListener) *ChatListenerManager {
	return &ChatListenerManager{
		repoChat:      repoChat,
		listener:      listener,
		chatListeners: make(map[int]*ChatListener),
	}
}

func (m *ChatListenerManager) GetChatListener(chatId int) (*ChatListener, error) {
	listener, exists := m.chatListeners[chatId]
	if exists {
		return listener, nil
	}

	if _, err := m.repoChat.GetChatByID(context.Background(), chatId); err != nil {
		return nil, err
	}

	// запустить прослушивание стораджа и отправлять в канал сообщений

	chatListener := NewChatListener(m.listener, chatId)
	go chatListener.StartListening(context.Background())

	m.chatListeners[chatId] = chatListener

	return chatListener, nil
}
