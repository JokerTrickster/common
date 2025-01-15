package kakao

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/JokerTrickster/common/oauth"
)

type KakaoService struct {
	appID    int64
	initOnce sync.Once
}

var kakaoInstance *KakaoService
var kakaoOnce sync.Once

// GetKakaoService returns the singleton instance of KakaoService
func GetKakaoService() *KakaoService {
	kakaoOnce.Do(func() {
		kakaoInstance = &KakaoService{}
	})
	return kakaoInstance
}

// Initialize initializes the Kakao OAuth configuration
func (s *KakaoService) Initialize(appID string) error {
	var err error
	s.initOnce.Do(func() {
		id, parseErr := strconv.ParseInt(appID, 10, 64)
		if parseErr != nil {
			err = fmt.Errorf("failed to parse Kakao app ID: %w", parseErr)
			return
		}
		s.appID = id
	})
	return err
}

// Validate validates the Kakao token and returns user data
func (s *KakaoService) Validate(ctx context.Context, token string) (oauth.OAuthData, error) {
	url := "https://kapi.kakao.com/v1/user/access_token_info"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return oauth.OAuthData{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return oauth.OAuthData{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return oauth.OAuthData{}, fmt.Errorf("invalid token: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return oauth.OAuthData{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var data struct {
		AppID int64 `json:"app_id"`
		ID    int64 `json:"id"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return oauth.OAuthData{}, fmt.Errorf("failed to parse response: %w", err)
	}

	if data.AppID != s.appID {
		return oauth.OAuthData{}, fmt.Errorf("invalid app ID: %d", data.AppID)
	}

	return oauth.OAuthData{
		ID:       strconv.FormatInt(data.ID, 10),
		Provider: "kakao",
	}, nil
}
