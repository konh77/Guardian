package token

import (
	"github.com/gin-gonic/gin"
)

type Request struct {
	GrantType    string `json:"grant_type" form:"grant_type"`
	Code         string `json:"code" form:"code"`
	RedirectURI  string `json:"redirect_uri" form:"redirect_uri"`
	ClientID     string `json:"client_id" form:"client_id"`
	ClientSecret string `json:"client_secret" form:"client_secret"`
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
	Scope        string `json:"scope" form:"scope"`
	Audience     string `json:"audience" form:"audience"`
	CodeVerifier string `json:"code_verifier" form:"code_verifier"`
	AuthMethod   string `json:"auth_method" form:"auth_method"`
}

func NewRequest(c *gin.Context) *Request {
	var req Request

	c.ShouldBind(&req)

	if req.GrantType == "" {
		req.GrantType = c.PostForm("grant_type")
	}
	if req.Code == "" {
		req.Code = c.PostForm("code")
	}
	if req.RedirectURI == "" {
		req.RedirectURI = c.PostForm("redirect_uri")
	}
	if req.RefreshToken == "" {
		req.RefreshToken = c.PostForm("refresh_token")
	}
	if req.Scope == "" {
		req.Scope = c.PostForm("scope")
	}
	if req.Audience == "" {
		req.Audience = c.PostForm("audience")
	}
	if req.CodeVerifier == "" {
		req.CodeVerifier = c.PostForm("code_verifier")
	}
	if req.AuthMethod == "" {
		req.AuthMethod = c.PostForm("auth_method")
	}
	if req.ClientID == "" {
		req.ClientID = c.PostForm("client_id")
	}

	return &req
}
