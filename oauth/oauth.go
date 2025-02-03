package oauth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
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

var httpClient = http.DefaultClient
var authProviderName = map[AuthProvider]string{
	AuthProviderGoogle: "google",
	AuthProviderKakao:  "kakao",
	AuthProviderNaver:  "naver",
}

// JwtVerify verify data
func JwtVerifyWithKeySet(ctx context.Context, p string, tokenString string, keySetUrl string) (jwt.MapClaims, error) {

	ctxHttp, ctxHttpCancel := context.WithTimeout(ctx, time.Second*10)
	defer ctxHttpCancel()

	req, err := http.NewRequestWithContext(ctxHttp, http.MethodGet, keySetUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	set, err := jwk.Parse(bytes)
	if err != nil {
		return nil, err
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify the token signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Retrieve the key ID from the token header
		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing key ID (kid) in token header")
		}

		// Look up the key in the key set
		key, exists := set.LookupKeyID(keyID)
		if !exists {
			return nil, fmt.Errorf("key ID %s not found", keyID)
		}

		var pubKey interface{}
		if err := key.Raw(&pubKey); err != nil {
			return nil, fmt.Errorf("failed to get raw key: %w", err)
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(token.Valid)
	return claims, nil
}
