package app

import (
	"chat-project/config"
	"chat-project/internal/controllers/restapi"
	"chat-project/internal/controllers/sse"
	"chat-project/internal/storage/postgres"
	redisStorage "chat-project/internal/storage/redis"
	"chat-project/internal/services"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Run(cfg *config.Config) {
	r := gin.Default()
	poolConfig, err := pgxpool.ParseConfig(cfg.Postgres.Url)
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_URL:", err)
	}

	pgPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}
	defer pgPool.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()

	chatListener := redisStorage.NewListener(redisClient)
	chatRepo := postgres.NewChatRepoPostgres(pgPool)

	service := services.New(chatRepo, chatListener)
	chatManager := services.NewChatListenerManager(chatRepo, chatListener)

	restapi.NewRouter(r, service)
	sse.NewRouter(r, chatManager)

	r.Run(fmt.Sprintf("%s:%s", "0.0.0.0", cfg.HTTP.Port))
}
