package authorize

import (
	"context"
	"net/http"

	"Guardian_/internal/oidc/authorize"

	"github.com/gin-gonic/gin"
)

type authorizeService interface {
	Authorize(ctx context.Context, req *authorize.AuthorizeRequest) (*authorize.Result, authorize.ErrorType, string, error)
}

type Handler struct {
	service authorizeService
}

func NewHandler(service authorizeService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Authorize(c *gin.Context) {
	req := &authorize.AuthorizeRequest{
		ClientID:            c.Query("client_id"),
		RedirectURI:         c.Query("redirect_uri"),
		ResponseType:        c.Query("response_type"),
		Scope:               c.Query("scope"),
		State:               c.Query("state"),
		Nonce:               c.Query("nonce"),
		CodeChallenge:       c.Query("code_challenge"),
		CodeChallengeMethod: c.Query("code_challenge_method"),
		UserID:              c.GetHeader("X-Authenticated-User"),
	}

	result, errType, description, err := h.service.Authorize(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "server_error"})
		return
	}
	if errType != authorize.ErrorNone {
		body := gin.H{"error": errType.String()}
		if description != "" {
			body["error_description"] = description
		}
		c.JSON(http.StatusBadRequest, body)
		return
	}

	c.Redirect(http.StatusFound, result.RedirectURL)
}
