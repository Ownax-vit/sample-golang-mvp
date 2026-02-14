package services

import (
	"context"
	"log"

	"chat-project/internal/domain"
	"chat-project/internal/storage"
)

type ClientConn chan domain.Message

// Прослушиватель сообщений в чате, который будет отправлять новые сообщения всем подписанным клиентам
type ChatListener struct {
	ChatId   int
	listener storage.ChatListener

	messages chan domain.Message

	// New client connections
	newClients chan ClientConn

	// Closed client connections
	closedClients chan ClientConn

	// Total client connections
	totalClients map[ClientConn]bool
	chatManager  *ChatListenerManager
}

func NewChatListener(listener storage.ChatListener, chatManager *ChatListenerManager, chatId int) *ChatListener {
	chatListener := &ChatListener{
		ChatId:        chatId,
		listener:      listener,
		messages:      make(chan domain.Message),
		newClients:    make(chan ClientConn),
		closedClients: make(chan ClientConn),
		totalClients:  make(map[ClientConn]bool),
		chatManager: chatManager,
	}

	return chatListener
}

func (l *ChatListener) StartListening(ctx context.Context) {
	go l.ListenChannels(ctx)
	go l.ListenStorage(ctx)
}

func (l *ChatListener) AddClient(ctx context.Context, clientChan ClientConn) {
	l.newClients <- clientChan
}

func (l *ChatListener) RemoveClient(ctx context.Context, clientChan ClientConn) {
	l.closedClients <- clientChan
}

func (l *ChatListener) BroadcastMessage(ctx context.Context, message domain.Message) {
	l.messages <- message
}

// Прослушивание новых сообщений из стораджа
func (l *ChatListener) ListenStorage(ctx context.Context) {
	messageChan := l.listener.Subscribe(ctx, l.ChatId)

	for {
		select {
		case msg := <-messageChan:
			l.BroadcastMessage(ctx, msg)
		case <-ctx.Done():
			log.Printf("Stopping storage listener for chat %d", l.ChatId)
			return
		}
	}
}

// Прослушивание каналов для управления клиентами и рассылки сообщений
func (l *ChatListener) ListenChannels(ctx context.Context) (<-chan domain.Message, error) {
	for {
		select {
		// Add new available client
		case client := <-l.newClients:
			l.totalClients[client] = true
			log.Printf("Client added to chat %d. %d registered clients", l.ChatId, len(l.totalClients))

		// Remove closed client
		case client := <-l.closedClients:
			delete(l.totalClients, client)
			close(client)
			log.Printf("Removed client from chatId %d. %d registered clients", l.ChatId, len(l.totalClients))

			if len(l.totalClients) == 0 {
				err := l.chatManager.CloseChat(l.ChatId)
				if err != nil {
					log.Printf("Error while closing chat %d: %v", l.ChatId, err)
				}
				return nil, nil
			}

		// Broadcast message to client
		case eventMsg := <-l.messages:
			for clientMessageChan := range l.totalClients {
				select {
				case clientMessageChan <- eventMsg:
					log.Printf("Message sent to client")
					// Message sent successfully
				default:
					log.Println("Failed to send message to client")
					// Failed to send, dropping message
				}
			}
		}
	}
}
