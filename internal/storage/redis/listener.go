package redis

import (
	"chat-project/internal/domain"
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type ListenerRedis struct {
	client *redis.Client
}

func NewListener(client *redis.Client) ListenerRedis {
	return ListenerRedis{
		client: client,
	}
}

func (l ListenerRedis) Subscribe(ctx context.Context, chatId int) <-chan domain.Message {
	chatStr := fmt.Sprintf("%d", chatId)
	pubsub := l.client.Subscribe(ctx, chatStr)

	ch := make(chan domain.Message)

	go func() {
		for {
			newMsg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				panic(err)
			}

			var msg domain.Message
			err = json.Unmarshal([]byte(newMsg.Payload), &msg)
			if err != nil {
				fmt.Printf("Error while unmarshall msg %v from channel %d", newMsg.Payload, chatId)
				continue
			}

			ch <- msg

		}
	}()

	return ch
}

func (l ListenerRedis) Publish(ctx context.Context, chatId int, message domain.Message) error {
	chatStr := fmt.Sprintf("%d", chatId)
	decodedMsg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = l.client.Publish(ctx, chatStr, decodedMsg).Err()
	if err != nil {
		return err
	}

	return nil
}
