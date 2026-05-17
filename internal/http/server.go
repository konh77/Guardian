package http

import (
	"Guardian_/internal/oidc/token"
	"Guardian_/internal/oidc/wellknown"
	"Guardian_/internal/store/memory"
	"Guardian_/internal/userinfo/source"
	"net/http"

	httpauthorize "Guardian_/internal/http/authorize"
	httpjwks "Guardian_/internal/http/jwks"
	httptoken "Guardian_/internal/http/token"
	httpwellknown "Guardian_/internal/http/wellknown"
	"Guardian_/internal/oidc/authorize"
	"Guardian_/internal/oidc/jwks"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine

	// Stores (exposed for testing/setup)
	ClientStore  *memory.ClientStore
	UserProvider *source.MemoryProvider
}

func NewServer(issuer string) *Server {
	r := gin.Default()

	// Initialize stores
	clientStore := memory.NewClientStore()
	userProvider := source.NewMemoryProvider()

	s := &Server{
		router:       r,
		ClientStore:  clientStore,
		UserProvider: userProvider,
	}

	s.setupRoutes(issuer)
	return s
}

func (s *Server) setupRoutes(issuer string) {
	s.router.GET("/hello", s.handleHello)
	s.router.GET("/user/:name", s.handleUser)
	s.router.GET("/teapot", s.teapot)

	wellknownService := wellknown.NewService("http://localhost:8080")
	wellknownHandler := httpwellknown.NewHandler(wellknownService)

	jwksService := jwks.NewService()
	jwksHandler := httpjwks.NewHandler(jwksService)
	s.router.GET("/jwks.json", jwksHandler.GetKeySet)

	authorizeService := authorize.NewService()
	authorizeHandler := httpauthorize.NewHandler(authorizeService)
	s.router.GET("/authorize", authorizeHandler.Authorize)
	tokenService := token.NewService()
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
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>This site can’t be reached</title>
<style>
	body {
		font-family: Arial, sans-serif;
		background-color: #f2f2f2;
		color: #202124;
		text-align: center;
		margin-top: 100px;
	}
	.container {
		max-width: 600px;
		margin: auto;
	}
	h1 {
		font-size: 24px;
		font-weight: normal;
	}
	p {
		color: #5f6368;
	}
	.code {
		margin-top: 20px;
		color: #5f6368;
		font-size: 14px;
	}
	button {
		margin-top: 20px;
		padding: 10px 20px;
		background-color: #1a73e8;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
	}
	button:hover {
		background-color: #1669c1;
	}
</style>
</head>
<body>
	<div class="container">
		<h1>This site can’t be reached</h1>
		<p>Check if there is a typo in example.jp.</p>
		<div class="code">DNS_PROBE_FINISHED_NXDOMAIN</div>
		<button onclick="location.reload()">Reload</button>
	</div>
</body>
</html>
`
	c.Data(404, "text/html; charset=utf-8", []byte(html))
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
