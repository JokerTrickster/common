package oauth

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// OAuthData represents user information returned by an OAuth provider
type OAuthData struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
}

// AuthProvider represents the available authentication providers
type AuthProvider uint8

const (
	AuthProviderGoogle AuthProvider = iota
	AuthProviderKakao
	AuthProviderNaver
)

var authProviderName = map[AuthProvider]string{
	AuthProviderGoogle: "google",
	AuthProviderKakao:  "kakao",
	AuthProviderNaver:  "naver",
}

// jwtVerifyWithKeySet verifies the JWT token using the given provider's public key URL
func JwtVerifyWithKeySet(ctx context.Context, provider string, token, keyURL string) (jwt.MapClaims, error) {
	// Fetch the key set from the provided URL
	keySet, err := fetchKeySet(keyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch key set for provider %s: %w", provider, err)
	}

	// Parse the JWT token
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if kid, ok := t.Header["kid"].(string); ok {
			if key, exists := keySet[kid]; exists {
				return key, nil
			}
		}
		return nil, errors.New("invalid key ID in token")
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT for provider %s: %w", provider, err)
	}

	// Extract claims if the token is valid
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// fetchKeySet fetches the public key set from the provided URL
func fetchKeySet(keyURL string) (map[string]interface{}, error) {
	// This function should fetch the key set from the given key URL
	// and return a map of keys indexed by their `kid`.
	// This is a placeholder implementation and should be replaced
	// with actual logic to fetch and parse the key set.

	// Example:
	// {
	//   "keys": [
	//     { "kid": "key1", "n": "...", "e": "..." },
	//     { "kid": "key2", "n": "...", "e": "..." }
	//   ]
	// }

	// Return an empty map as a placeholder
	return map[string]interface{}{
		"example-kid": "example-key", // Replace with actual key fetching logic
	}, nil
}
