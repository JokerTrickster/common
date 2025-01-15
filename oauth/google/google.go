package google

import (
	"context"
	"fmt"
	"sync"

	"github.com/JokerTrickster/common/oauth"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleService struct {
	config   *oauth2.Config
	authMeta AuthMeta
	initOnce sync.Once
}

var googleInstance *GoogleService
var googleOnce sync.Once

type AuthMeta struct {
	GoogleIosID string
	GoogleAndID []string
}

// GetGoogleService returns the singleton instance of GoogleService
func GetGoogleService() *GoogleService {
	googleOnce.Do(func() {
		googleInstance = &GoogleService{}
	})
	return googleInstance
}

// Initialize initializes the Google OAuth configuration
func (s *GoogleService) Initialize(clientID, clientSecret, redirectURL string, googleIosID string, googleAndIDs []string) {
	s.initOnce.Do(func() {
		s.config = &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		}
		s.authMeta = AuthMeta{
			GoogleIosID: googleIosID,
			GoogleAndID: googleAndIDs,
		}
	})
}

// Validate validates the Google token and returns user data
func (s *GoogleService) Validate(ctx context.Context, token string) (oauth.OAuthData, error) {
	claims, err := oauth.JwtVerifyWithKeySet(ctx, "google", token, "https://www.googleapis.com/oauth2/v3/certs")
	if err != nil {
		return oauth.OAuthData{}, err
	}

	aud, okAud := claims["aud"].(string)
	iss, okIss := claims["iss"].(string)
	sub, okSub := claims["sub"].(string)
	email, okEmail := claims["email"].(string)

	if !okAud || !okIss || !okSub || !okEmail ||
		(aud != s.authMeta.GoogleIosID && s.isGoogleIDNotExisted(aud)) ||
		(iss != "accounts.google.com" && iss != "https://accounts.google.com") {
		return oauth.OAuthData{}, fmt.Errorf("invalid token claims: %+v", claims)
	}

	return oauth.OAuthData{
		ID:       sub,
		Email:    email,
		Provider: "google",
	}, nil
}

func (s *GoogleService) isGoogleIDNotExisted(key string) bool {
	for _, id := range s.authMeta.GoogleAndID {
		if key == id {
			return false
		}
	}
	return true
}

// ExchangeToken exchanges an authorization code for an access token
func (s *GoogleService) ExchangeToken(ctx context.Context, authCode string) (*oauth2.Token, error) {
	// GoogleConfig가 초기화되어 있는지 확인
	if s.config == nil {
		return nil, fmt.Errorf("Google OAuth configuration is not initialized")
	}

	// Authorization Code로 Access Token 교환
	token, err := s.config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	return token, nil
}
