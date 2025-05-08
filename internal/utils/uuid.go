package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// SafeUUIDString returns a string representation of a UUID pointer, handling nil cases
func SafeUUIDString(id *uuid.UUID) string {
	if id == nil {
		return "nil"
	}
	return id.String()
}

// SafeString returns a string representation of a string pointer, handling nil cases
func SafeString(s *string) string {
	if s == nil {
		return "nil"
	}
	return *s
}

// NullUUIDToPtr converts a nullable pgtype.UUID to a *uuid.UUID
func NullUUIDToPtr(id pgtype.UUID) *uuid.UUID {
	if !id.Valid {
		return nil
	}
	result := uuid.UUID(id.Bytes)
	return &result
}
