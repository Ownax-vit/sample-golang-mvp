package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"chat-project/internal/domain"
	"chat-project/internal/storage"
)

type ChatRepoPostgres struct {
	pool *pgxpool.Pool
}

func NewChatRepoPostgres(pgpool *pgxpool.Pool) *ChatRepoPostgres {
	return &ChatRepoPostgres{
		pool: pgpool,
	}
}

func (r ChatRepoPostgres) CreateChat(ctx context.Context, chat domain.Chat) (domain.Chat, error) {
	_, err := r.pool.Exec(ctx, "INSERT INTO chats (id, title, created_at) VALUES ($1, $2, $3)", chat.ID, chat.Title, chat.CreatedAt)
	if err != nil {
		return domain.Chat{}, err
	}

	return chat, nil
}

func (r ChatRepoPostgres) GetChatByID(ctx context.Context, chatId int) (domain.Chat, error) {
	var chat domain.Chat
	err := r.pool.QueryRow(ctx, "SELECT id, title, created_at FROM chats WHERE id = $1", chatId).Scan(
		&chat.ID, &chat.Title, &chat.CreatedAt,
	)
	if err != nil {
		return domain.Chat{}, fmt.Errorf("error while getting chat: %s %w", err, storage.ChatNotFoundError)
	}

	return chat, nil
}

func (r ChatRepoPostgres) AddMessage(ctx context.Context, message domain.Message, chatId int) (domain.Message, error) {
	_, err := r.pool.Exec(
		ctx,
		"INSERT INTO messages (id, chat_id, text) VALUES ($1, $2, $3)", message.ID, chatId, message.Text,
	)
	if err != nil {
		return domain.Message{}, err
	}

	return message, nil
}

func (r ChatRepoPostgres) GetWithMessages(ctx context.Context, chatId int) (domain.Chat, error) {
	var chat domain.Chat
	query := `
	SELECT
	c.id,
	c.title,
	c.created_at,
	COALESCE(
		jsonb_agg(
			jsonb_build_object(
				'id', m.id,
				'text', m.text,
				'created_at',  m.created_at
			)
		) FILTER (WHERE m.id IS NOT NULL), '[]') AS messages
	FROM chats c
	LEFT JOIN messages m ON c.id = m.chat_id
	WHERE c.id = $1
	GROUP BY c.id, c.created_at, c.title
	`
	err := r.pool.QueryRow(ctx, query, chatId).Scan(&chat.ID, &chat.Title, &chat.CreatedAt, &chat.Messages)
	if err != nil {
		return domain.Chat{}, fmt.Errorf("error while getting chat with messages: %s %w", err, storage.ChatNotFoundError)
	}

	return chat, nil
}

func (r ChatRepoPostgres) DeleteChat(ctx context.Context, chatId int) error {
	//
	_, err := r.pool.Exec(ctx, "DELETE FROM chats WHERE id = $1", chatId)
	return err
}
