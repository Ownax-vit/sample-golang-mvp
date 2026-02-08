package app

import (
	"chat-project/config"
	"chat-project/internal/controllers/restapi"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	// Application entry point
	r := gin.Default()

	restapi.NewRouter(r, cfg)

	r.Run(fmt.Sprintf("%s:%s", "0.0.0.0", cfg.HTTP.Port))
}
