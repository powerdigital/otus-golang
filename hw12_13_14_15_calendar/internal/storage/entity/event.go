package entity

import "time"

type Event struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"userId" db:"user_id"`
	Title        string    `json:"title" db:"title"`
	EventTime    time.Time `json:"eventTime" db:"event_time"`
	DurationMin  int       `json:"durationMin" db:"duration_min"`
	NoticeBefore int       `json:"noticeBefore" db:"notice_before"`
	Description  string    `json:"description" db:"description"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}
