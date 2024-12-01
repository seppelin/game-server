package models

import "time"

type SessionID int64

type Session struct {
	Id            SessionID
	Name          string
	UserID        UserID
	CreateAddress string
	CreatedAt     time.Time
	AccessedAt    time.Time
	ExpiresAt     time.Time
}
