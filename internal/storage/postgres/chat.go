package postgres

import (
	"chat-project/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
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
	_, err := r.pool.Exec(ctx, "INSERT INTO chats (id, name) VALUES ($1, $2)", chat.ID, chat.Title)
	if err != nil {
		return domain.Chat{}, err
	}

	return chat, nil
}

func (r ChatRepoPostgres) GetChatByID(ctx context.Context, chatId int) (domain.Chat, error) {
	var chat domain.Chat
	err := r.pool.QueryRow(ctx, "SELECT id, name FROM chats WHERE id = $1", chatId).Scan(&chat.ID, &chat.Title)
	if err != nil {
		return domain.Chat{}, err
	}

	return chat, nil
}

func (r ChatRepoPostgres) AddMessage(ctx context.Context, message domain.Message, chatId int) (domain.Message, error) {
	_, err := r.pool.Exec(
		ctx,
		"INSERT INTO messages (id, chat_id, content) VALUES ($1, $2, $3)", message.ID, chatId, message.Text,
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
	c.name,
	COALESCE(json_agg(json_build_object(m.id, 'c.id', 'm.text', m.created_at)), '[]') AS messages
	FROM chats c
	LEFT JOIN messages m ON c.id = m.chat_id
	WHERE c.id = $1
	`
	err := r.pool.QueryRow(ctx, query, chatId).Scan(&chat.ID, &chat.Title, &chat.Messages)
	if err != nil {
		return domain.Chat{}, err
	}

	return chat, nil
}

func (r ChatRepoPostgres) DeleteChat(ctx context.Context, chatId int) error {
	//
	_, err := r.pool.Exec(ctx, "DELETE FROM chats WHERE id = $1", chatId)
	return err
}


