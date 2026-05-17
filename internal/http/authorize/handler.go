package authorize

import (
	"Guardian_/internal/oidc/authorize"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *authorize.Service
}

func NewHandler(service *authorize.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Authorize(c *gin.Context) {
	result, err := h.service.Authorize()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(302, result.RedirectURL)
}
