package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func SafeUUIDString(id *uuid.UUID) string {
	if id == nil {
		return "nil"
	}
	return id.String()
}

func SafeString(s *string) string {
	if s == nil {
		return "nil"
	}
	return *s
}

func NullUUIDToPtr(id pgtype.UUID) *uuid.UUID {
	if !id.Valid {
		return nil
	}
	result := uuid.UUID(id.Bytes)
	return &result
}
