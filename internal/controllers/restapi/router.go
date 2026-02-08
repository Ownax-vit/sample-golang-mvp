package restapi

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"chat-project/config"
	v1 "chat-project/internal/controllers/restapi/v1"
	"chat-project/internal/services"
	"chat-project/internal/storage/postgres"

	docs "chat-project/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter создает и настраивает роутер для API
func NewRouter(app *gin.Engine, cfg *config.Config) {

	poolConfig, err := pgxpool.ParseConfig(cfg.Postgres.Url)
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_URL:", err)
	}

	pgPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}

	chatRepo := postgres.NewChatRepoPostgres(pgPool)
	chatController := v1.NewChatController(services.New(chatRepo))
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
