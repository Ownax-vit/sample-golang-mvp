package app

import (
	"chat-project/config"
	"chat-project/internal/controllers/restapi"
	"chat-project/internal/controllers/sse"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	r := gin.Default()

	restapi.NewRouter(r, cfg)
	sse.NewRouter(r, cfg)

	r.Run(fmt.Sprintf("%s:%s", "0.0.0.0", cfg.HTTP.Port))
}
