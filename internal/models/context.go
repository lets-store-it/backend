package models

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// Context keys
	OrganizationIDContextKey ContextKey = "organization-id"
	UserIDContextKey         ContextKey = "user-id"
	IsSystemUserContextKey   ContextKey = "is-system-user"
	SetCookieContextKey      ContextKey = "set-cookie"
)
