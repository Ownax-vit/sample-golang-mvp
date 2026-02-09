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
	ChatId int
	listener     storage.ChatListener

	NewMessages chan domain.Message

	// New client connections
	NewClients chan ClientConn

	// Closed client connections
	ClosedClients chan ClientConn

	// Total client connections
	TotalClients map[ClientConn]bool
}

func NewChatListener(listener storage.ChatListener, chatId int) *ChatListener{
	chatListener := &ChatListener{
		ChatId:        chatId,
		listener:            listener,
		NewMessages:   make(chan domain.Message),
		NewClients:    make(chan ClientConn),
		ClosedClients: make(chan ClientConn),
		TotalClients:  make(map[ClientConn]bool),
	}

	return chatListener
}

func (l *ChatListener) StartListening(ctx context.Context) {
	go l.ListenChannels(ctx)
	go l.ListenStorage(ctx)
}

func (l *ChatListener) AddClient(ctx context.Context, clientChan ClientConn) {
	l.NewClients <- clientChan
}

func (l *ChatListener) RemoveClient(ctx context.Context, clientChan ClientConn) {
	l.ClosedClients <- clientChan
}

func (l *ChatListener) BroadcastMessage(ctx context.Context, message domain.Message) {
	l.NewMessages <- message
}

// Прослушивание новых сообщений из стораджа
func (l *ChatListener) ListenStorage(ctx context.Context){
	messageChan, err := l.listener.ListenChat(ctx, l.ChatId)

	if err != nil {
		log.Printf("Error listening storage for chat %d: %v", l.ChatId, err)
		return
	}

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
		case client := <-l.NewClients:
			l.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(l.TotalClients))

		// Remove closed client
		case client := <-l.ClosedClients:
			delete(l.TotalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(l.TotalClients))

		// Broadcast message to client
		case eventMsg := <-l.NewMessages:
			for clientMessageChan := range l.TotalClients {
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
