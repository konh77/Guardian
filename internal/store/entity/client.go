package entity

import "errors"

var ErrClientNotFound = errors.New("client not found")

type Client struct {
	ID                      string   `json:"client_id"`
	Secret                  string   `json:"client_secret,omitempty"`
	RedirectURIs            []string `json:"redirect_uris"`
	Name                    string   `json:"client_name,omitempty"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method,omitempty"`
	ClientType              string   `json:"client_type,omitempty"`
}

func (c *Client) IsPublicClient() bool {
	return c.ClientType == "public"
}

func (c *Client) IsConfidentialClient() bool {
	return c.ClientType == "confidential" || c.ClientType == ""
}
