package models

type ContextKey string

const (
	OrganizationIDContextKey ContextKey = "organization-id"
	UserIDContextKey         ContextKey = "user-id"
	TvBoardIDContextKey      ContextKey = "tv-board-id"
	IsSystemUserContextKey   ContextKey = "is-system-user"
	SetCookieContextKey      ContextKey = "set-cookie"
)
