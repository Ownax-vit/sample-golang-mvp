package restapi

import (
	"github.com/gin-gonic/gin"

	docs "chat-project/docs"
	v1 "chat-project/internal/controllers/restapi/v1"
	"chat-project/internal/services"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter создает и настраивает роутер для API
func NewRouter(app *gin.Engine, service *services.ChatService) {
	chatController := v1.NewChatController(service)
	docs.SwaggerInfo.BasePath = "/v1"
	// Routers
	apiV1Group := app.Group("/v1")
	{
		chats := apiV1Group.Group("/chats")
		{
			chats.POST("/", chatController.CreateChat)
			chats.GET("/:chatId", chatController.GetChat)
			chats.POST("/:chatId/messages", chatController.AddMessage)
			chats.DELETE("/:chatId", chatController.DeleteChat)
		}
	}

	apiV1Group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}
