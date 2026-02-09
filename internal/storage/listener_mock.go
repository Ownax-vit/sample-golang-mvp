package storage

import (
	"chat-project/internal/domain"
	"context"
	"time"
)

type ListenerMock struct {
}

func NewListenerMock() *ListenerMock {
	return &ListenerMock{}
}

func (l ListenerMock) ListenChat(ctx context.Context, chatId int) (<-chan domain.Message, error) {
	// использовать subsribe на канал сообщений, который будет отправлять новые сообщения в чат
	ch := make(chan domain.Message)

	// TODO  LISTEN/NOTIFY в PostgreSQL, PUB Sub REDIS или другой механизм уведомлений.

	go func() {
		for {
			ch <- domain.Message{
				ID:   0,
				Text: "New message in chat",
			}
			time.Sleep(5 * time.Second) // имитация получения новых сообщений каждые 5 секунд
		}
	}()

	return ch, nil
}
