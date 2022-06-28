package entity

type Event struct {
	ID           int64  `json:"id" db:"id"`
	UserID       int64  `json:"userId" db:"user_id"`
	Title        string `json:"title" db:"title"`
	EventTime    string `json:"eventTime" db:"event_time"`
	DurationMin  int64  `json:"durationMin" db:"duration_min"`
	NoticeBefore int64  `json:"noticeBefore" db:"notice_before"`
	Description  string `json:"description" db:"description"`
	CreatedAt    string `json:"createdAt" db:"created_at"`
}
