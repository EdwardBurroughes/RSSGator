// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Feed struct {
	ID            int32
	Name          string
	Url           string
	UserID        uuid.NullUUID
	LastFetchedAt sql.NullTime
}

type FeedFollow struct {
	ID        int32
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    int32
}

type Post struct {
	ID          int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       sql.NullString
	Url         sql.NullString
	Description sql.NullString
	PublishedAt sql.NullString
	FeedID      sql.NullInt32
}

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}
