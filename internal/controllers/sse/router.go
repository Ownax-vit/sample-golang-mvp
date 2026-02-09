package sse

import (
	"context"
	"io"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"chat-project/config"
	"chat-project/internal/services"
	"chat-project/internal/storage"
	"chat-project/internal/storage/postgres"
)

var chatManager *services.ChatListenerManager

func NewRouter(app *gin.Engine, cfg *config.Config) {
	router := app.Group("/sse")
	poolConfig, err := pgxpool.ParseConfig(cfg.Postgres.Url)
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_URL:", err)
	}

	pgPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}

	chatListener := storage.NewListenerMock()
	chatRepo := postgres.NewChatRepoPostgres(pgPool)
	chatManager = services.NewChatListenerManager(chatRepo, chatListener)

	router.GET("/sse", HeadersMiddleware(), serveHTTP(), func(c *gin.Context) {
		v, ok := c.Get("clientChan")
		if !ok {
			return
		}
		clientChan, ok := v.(services.ClientConn)
		if !ok {
			return
		}
		c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			if msg, ok := <-clientChan; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})
}

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}

func serveHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize client channel

		chatIdStr := c.Query("chatId")
		if chatIdStr == "" {
			c.AbortWithStatus(400)
			return
		}

		chatId, err := strconv.Atoi(chatIdStr)
		if err != nil {
			c.AbortWithStatus(400)
			return
		}

		clientChan := make(services.ClientConn)

		// Send new connection to event server
		chatListener, err := chatManager.GetChatListener(chatId)
		if err != nil {
			c.AbortWithStatus(404)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		chatListener.AddClient(ctx, clientChan)

		go func() {
			<-c.Request.Context().Done()
			cancel()
			chatListener.RemoveClient(c, clientChan)
		}()

		c.Set("clientChan", clientChan)

		c.Next()
	}
}
