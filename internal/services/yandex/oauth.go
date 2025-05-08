package yandex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrInvalidOrExpiredToken = errors.New("token invalid or expired")
)

type YandexOAuthService struct {
	clientID     string
	clientSecret string
	tracer       trace.Tracer
}

type YandexOAuthUserInfo struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	DefaultEmail string `json:"default_email"`
}

type YandexOAuthServiceConfig struct {
	ClientID     string
	ClientSecret string
}

func NewYandexOAuthService(config YandexOAuthServiceConfig) *YandexOAuthService {
	return &YandexOAuthService{
		clientID:     config.ClientID,
		clientSecret: config.ClientSecret,
		tracer:       otel.GetTracerProvider().Tracer("yandex-oauth"),
	}
}

func (s *YandexOAuthService) GetUserInfo(ctx context.Context, accessToken string) (YandexOAuthUserInfo, error) {
	ctx, span := s.tracer.Start(ctx, "GetUserInfo",
		trace.WithAttributes(
			attribute.String("client_id", s.clientID),
		),
	)
	defer span.End()

	url := fmt.Sprintf("https://login.yandex.ru/info?format=json&jwt_secret=%s", s.clientSecret)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create request")
		return YandexOAuthUserInfo{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to execute request")
		return YandexOAuthUserInfo{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		span.SetStatus(codes.Error, "token invalid or expired")
		return YandexOAuthUserInfo{}, ErrInvalidOrExpiredToken
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		span.RecordError(err)
		span.SetStatus(codes.Error, "unexpected status code")
		return YandexOAuthUserInfo{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to read response body")
		return YandexOAuthUserInfo{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo YandexOAuthUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to unmarshal response body")
		return YandexOAuthUserInfo{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	span.SetAttributes(
		attribute.String("user.id", userInfo.ID),
		attribute.String("user.email", userInfo.DefaultEmail),
	)
	span.SetStatus(codes.Ok, "successfully retrieved user info")

	return userInfo, nil
}
