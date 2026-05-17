package main

import (
	"Guardian_/internal/store/entity"
	"log"

	"Guardian_/internal/http"
)

func main() {
	issuer := "http://localhost:8080"
	server := http.NewServer(issuer)

	server.ClientStore.RegisterClient(&entity.Client{
		ID:           "my-client",
		Secret:       "my-secret",
		RedirectURIs: []string{"http://localhost:3000/callback"},
		Name:         "Test Client",
		ClientType:   "confidential",
	})

	server.UserProvider.RegisterUser("user123", map[string]any{
		"sub":   "user123",
		"name":  "Test User",
		"email": "test@example.com",
	})

	log.Println("Starting server on :8080")
	log.Println("Test client: my-client")
	log.Println("Test user: user123")

	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
