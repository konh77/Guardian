package http

import (
	"context"
	"net/http"

	httpauthorize "Guardian_/internal/http/authorize"
	httpjwks "Guardian_/internal/http/jwks"
	httptoken "Guardian_/internal/http/token"
	httpwellknown "Guardian_/internal/http/wellknown"
	"Guardian_/internal/oidc/authorize"
	"Guardian_/internal/oidc/jwks"
	"Guardian_/internal/oidc/token"
	"Guardian_/internal/oidc/wellknown"
	"Guardian_/internal/store"
	"Guardian_/internal/store/entity"
	"Guardian_/internal/userinfo"

	"github.com/gin-gonic/gin"
)

type clientSeeder interface {
	store.ClientStore
	SaveClient(ctx context.Context, client *entity.Client) error
}

type userSeeder interface {
	userinfo.Provider
	RegisterUser(ctx context.Context, sub string, claims map[string]any) error
}

type backend struct {
	clients   clientSeeder
	users     userSeeder
	authCodes store.AuthCodeStore
	nonces    store.NonceStore
}

type Server struct {
	router *gin.Engine

	// Stores (exposed for testing/setup)
	ClientStore  clientSeeder
	UserProvider userSeeder

	authCodeStore store.AuthCodeStore
	nonceStore    store.NonceStore
}

func NewServer(issuer string) *Server {
	r := gin.Default()

	b := newBackend()

	s := &Server{
		router:        r,
		ClientStore:   b.clients,
		UserProvider:  b.users,
		authCodeStore: b.authCodes,
		nonceStore:    b.nonces,
	}

	s.setupRoutes(issuer)
	return s
}

func (s *Server) setupRoutes(issuer string) {
	s.router.GET("/hello", s.handleHello)
	s.router.GET("/user/:name", s.handleUser)
	s.router.GET("/teapot", s.teapot)
	s.router.HEAD("/teapot", s.teapot)

	wellknownService := wellknown.NewService(issuer)
	wellknownHandler := httpwellknown.NewHandler(wellknownService)

	jwksService := jwks.NewService()
	jwksHandler := httpjwks.NewHandler(jwksService)
	s.router.GET("/jwks.json", jwksHandler.GetKeySet)

	authorizeService := authorize.NewService(
		s.ClientStore,
		s.UserProvider,
		s.authCodeStore,
		s.nonceStore,
	)
	authorizeHandler := httpauthorize.NewHandler(authorizeService)
	s.router.GET("/authorize", authorizeHandler.Authorize)
	tokenService := token.NewService(s.authCodeStore, jwksService, issuer)
	tokenHandler := httptoken.NewHandler(tokenService)
	s.router.POST("/token", tokenHandler.ExchangeToken)
	s.router.GET("/.well-known/openid-configuration", wellknownHandler.GetConfiguration)
}

func (s *Server) handleHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello!",
	})
}

func (s *Server) teapot(c *gin.Context) {
	c.JSON(http.StatusTeapot, gin.H{
		"error":   "I’m a teapot",
		"message": "The server refuses to brew coffee because it is, permanently, a teapot.",
		"status":  http.StatusTeapot,
	})
}

func (s *Server) handleUser(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome, " + name + "!",
	})
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
