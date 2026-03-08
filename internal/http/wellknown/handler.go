package jwks

import "github.com/gin-gonic/gin"
import "net/http"
import "Guardian_/internal/oidc/wellknown"

type Handler struct {
	service *wellknown.Service
}

func NewHandler(service *wellknown.Service) *Handler {
	return &Handler{
		service: service,
	}

}

func (h *Handler) GetKeySet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hi",
	})
}
