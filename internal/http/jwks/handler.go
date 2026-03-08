package jwks

import "github.com/gin-gonic/gin"
import "net/http"
import "Guardian_/internal/oidc/jwks"

type Handler struct {
	service *jwks.Service
}

func NewHandler(service *jwks.Service) *Handler {
	return &Handler{
		service: service,
	}

}

func (h *Handler) GetKeySet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hi",
	})
}
