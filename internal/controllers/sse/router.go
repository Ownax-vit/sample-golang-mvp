package sse

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"

	"chat-project/internal/services"
)

func NewRouter(app *gin.Engine, chatManager *services.ChatListenerManager) {
	router := app.Group("/sse")

	sseController := &SSEController{
		chatManager: chatManager,
	}

	router.GET("/sse", HeadersMiddleware(), sseController.serveHTTP(), func(c *gin.Context) {
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

type SSEController struct {
	chatManager *services.ChatListenerManager
}

func (sse *SSEController) serveHTTP() gin.HandlerFunc {
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
		chatListener, err := sse.chatManager.GetChatListener(chatId)
		if err != nil {
			c.AbortWithStatus(404)
			return
		}

		chatListener.AddClient(c, clientChan)

		go func() {
			<-c.Request.Context().Done()
			chatListener.RemoveClient(c, clientChan)
		}()

		c.Set("clientChan", clientChan)

		c.Next()
	}
}
