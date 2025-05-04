package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UuidFromPgx(id pgtype.UUID) *uuid.UUID {
	if !id.Valid {
		return nil
	}
	result := uuid.UUID(id.Bytes)
	return &result
}

func PgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func PgText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

func PgTextPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: *s != ""}
}

// NullablePgUUID converts a pointer to UUID to pgtype.UUID
func NullablePgUUID(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return PgUUID(*id)
}

// IsValidUUID checks if a UUID is not nil
func IsValidUUID(id uuid.UUID) bool {
	return id != uuid.Nil
}
