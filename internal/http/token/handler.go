package token

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Guardian_/internal/oidc/signing"
	"Guardian_/internal/oidc/token"
	"Guardian_/internal/store"
	"Guardian_/internal/userinfo"
)

// Handler handles HTTP requests for token endpoint
type Handler struct {
	service *token.Service
}

// NewHandler creates a new token handler
func NewHandler(
	issuer string,
	clientStore store.ClientStore,
	authCodeStore store.AuthCodeStore,
	tokenStore store.TokenStore,
	userinfoProvider userinfo.Provider,
	signingService *signing.Service,
) *Handler {
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
