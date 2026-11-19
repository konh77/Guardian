package main

import (
	"context"
	"log"

	"Guardian_/internal/http"
	"Guardian_/internal/store/entity"
)

func main() {
	issuer := "http://localhost:8080"
	server := http.NewServer(issuer)

	ctx := context.Background()

	if err := server.ClientStore.SaveClient(ctx, &entity.Client{
		ID:           "my-client",
		Secret:       "my-secret",
		RedirectURIs: []string{"http://localhost:3000/callback"},
		Name:         "Test Client",
		ClientType:   "confidential",
	}); err != nil {
		log.Fatalf("seed client: %v", err)
	}

	if err := server.UserProvider.RegisterUser(ctx, "user123", map[string]any{
		"sub":   "user123",
		"name":  "Test User",
		"email": "test@example.com",
	}); err != nil {
		log.Fatalf("seed user: %v", err)
	}

	log.Println("Starting server on :8080")
	log.Println("Test client: my-client")
	log.Println("Test user: user123")

	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
