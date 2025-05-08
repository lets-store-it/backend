package database

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// UUID
func UuidFromPgx(id pgtype.UUID) uuid.UUID {
	if !id.Valid {
		return uuid.Nil
	}
	return uuid.UUID(id.Bytes)
}

func UuidPtrFromPgx(id pgtype.UUID) *uuid.UUID {
	if !id.Valid {
		return nil
	}
	result := uuid.UUID(id.Bytes)
	return &result
}

func PgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func PgUuidPtr(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return PgUUID(*id)
}

// Text
func PgText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

func PgTextPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: *s != ""}
}
