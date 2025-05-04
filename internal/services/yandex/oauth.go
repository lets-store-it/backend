package yandex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrInvalidOrExpiredToken = errors.New("token invalid or expired")
)

type YandexOAuthService struct {
	clientID     string
	clientSecret string
}

type YandexOAuthUserInfo struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	DefaultEmail string `json:"default_email"`
}

func NewYandexOAuthService(clientID, clientSecret string) *YandexOAuthService {
	return &YandexOAuthService{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (s *YandexOAuthService) GetUserInfo(ctx context.Context, accessToken string) (*YandexOAuthUserInfo, error) {
	url := fmt.Sprintf("https://login.yandex.ru/info?format=json&jwt_secret=%s", s.clientSecret)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrInvalidOrExpiredToken
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo YandexOAuthUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &userInfo, nil
}
