package token

import (
	"net/http"

	"Guardian_/internal/oidc/token"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *token.Service
}

func NewHandler(service *token.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ExchangeToken(c *gin.Context) {
	code := c.PostForm("code")

	result, err := h.service.ExchangeToken(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
