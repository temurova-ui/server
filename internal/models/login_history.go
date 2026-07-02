package models

import (
	"time"
)

type LoginHistory struct {
	ID        int64
	UserID    int64
	IPAddress string
	UserAgent string
	CreatedAt time.Time
}
