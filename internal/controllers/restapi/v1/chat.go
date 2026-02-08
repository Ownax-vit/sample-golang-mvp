package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"chat-project/internal/dto"
	"chat-project/internal/services"
	"chat-project/internal/storage"
)

type ChatController struct {
	service *services.ChatService
}

func NewChatController(service *services.ChatService) *ChatController {
	return &ChatController{
		service: service,
	}
}

// CreateChat создает новый чат
//
//	@Summary      Создать чат
//	@Description  Создает новый чат с указанным названием
//	@Tags         chats
//	@Accept       json
//	@Produce      json
//	@Param        chat  body      dto.ChatIn  true  "Данные чата"
//	@Success      200   {object}  dto.ChatResponse
//	@Failure      400   {object}  map[string]string  "Неверный запрос"
//	@Failure      500   {object}  map[string]string  "Внутренняя ошибка сервера"
//	@Router       /chats [post]
func (c *ChatController) CreateChat(ctx *gin.Context) {
	var chat dto.ChatIn
	if err := ctx.ShouldBindJSON(&chat); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if chat.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	chatResponse, err := c.service.Create(
		ctx,
		chat,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, chatResponse)
}

// AddMessage добавляет сообщение в чат
//
//	@Summary      Добавить сообщение
//	@Description  Добавляет новое сообщение в указанный чат
//	@Tags         chats
//	@Accept       json
//	@Produce      json
//	@Param        chatId   path      int            true   "ID чата"
//	@Param        message  body      dto.MessageIn  true   "Данные сообщения"
//	@Success      200      {object}  dto.MessageResponse
//	@Failure      400      {object}  map[string]string  "Неверный запрос"
//	@Failure      500      {object}  map[string]string  "Внутренняя ошибка сервера"
//	@Router       /chats/{chatId}/messages [post]
func (c ChatController) AddMessage(ctx *gin.Context) {
	var message dto.MessageIn
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatIdParam := ctx.Param("chatId")
	chatId, err := dto.ParseID(chatIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	msg_resp, err := c.service.AddMessage(ctx, chatId, message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, msg_resp)

}

// GetChat получает чат с сообщениями
//
//	@Summary      Получить чат
//	@Description  Получает информацию о чате со всеми сообщениями
//	@Tags         chats
//	@Accept       json
//	@Produce      json
//	@Param        chatId  path      int  true  "ID чата"https://github.com/Ownax-vit
//	@Success      200     {object}  dto.ChatWithMessagesResponse
//	@Failure      400     {object}  map[string]string  "Неверный запрос"
//	@Failure      500     {object}  map[string]string  "Внутренняя ошибка сервера"
//	@Router       /chats/{chatId} [get]
func (c ChatController) GetChat(ctx *gin.Context) {
	chatIdParam := ctx.Param("chatId")
	chatId, err := dto.ParseID(chatIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	chatResp, err := c.service.GetWithMessages(ctx, chatId)
	if err != nil {
		if errors.Is(err, storage.ChatNotFoundError) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "chat not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, chatResp)
}

// DeleteChat удаляет чат
//
//	@Summary      Удалить чат
//	@Description  Удаляет чат по указанному ID
//	@Tags         chats
//	@Accept       json
//	@Produce      json
//	@Param        chatId  path      int  true  "ID чата"
//	@Success      204     "Чат успешно удален"
//	@Failure      400     {object}  map[string]string  "Неверный запрос"
//	@Failure      500     {object}  map[string]string  "Внутренняя ошибка сервера"
//	@Router       /chats/{chatId} [delete]
func (c ChatController) DeleteChat(ctx *gin.Context) {
	chatIdParam := ctx.Param("chatId")
	chatId, err := dto.ParseID(chatIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	err = c.service.DeleteChat(ctx, chatId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "chat deleted successfully"})
}
