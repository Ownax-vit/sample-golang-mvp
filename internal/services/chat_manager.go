package services

import (
	"chat-project/internal/storage"
	"context"
	"errors"
	"sync"
)

var ChatNotFoundError = errors.New("chat not found")

// Менеджер слушателей, который будет хранить всех слушателей для разных чатов и создавать новых при необходимости
type ChatListenerManager struct {
	repoChat       storage.ChatRepo
	listener       storage.ChatListener
	chatListeners  map[int]*ChatListener
	closedChannels chan int
	mu sync.Mutex
}

func NewChatListenerManager(repoChat storage.ChatRepo, listener storage.ChatListener) *ChatListenerManager {
	chatListener := &ChatListenerManager{
		repoChat:       repoChat,
		listener:       listener,
		chatListeners:  make(map[int]*ChatListener),
		closedChannels: make(chan int, 1),
	}
	go chatListener.ListenChannels()

	return chatListener
}

func (m *ChatListenerManager) ListenChannels() {
	for chatId := range m.closedChannels {
		m.mu.Lock()
		delete(m.chatListeners, chatId)
		m.mu.Unlock()
	}
}

func (m *ChatListenerManager) CloseChat(chatId int) error {
	m.mu.Lock()
	_, exists := m.chatListeners[chatId]
	m.mu.Unlock()
	if !exists {
		return ChatNotFoundError
	}

	m.closedChannels <- chatId
	return nil
}

func (m *ChatListenerManager) GetChatListener(chatId int) (*ChatListener, error) {
	m.mu.Lock()
	listener, exists := m.chatListeners[chatId]
	m.mu.Unlock()
	if exists {
		return listener, nil
	}

	if _, err := m.repoChat.GetChatByID(context.Background(), chatId); err != nil {
		return nil, err
	}

	// запустить прослушивание стораджа и отправлять в канал сообщений

	chatListener := NewChatListener(m.listener, m, chatId)

	m.mu.Lock()
	if _, exists := m.chatListeners[chatId]; exists {
		m.mu.Unlock()
		return m.chatListeners[chatId], nil
	}
	m.chatListeners[chatId] = chatListener
	m.mu.Unlock()

	go chatListener.StartListening(context.Background())
	return chatListener, nil
}
