package storage

import (
	"reflect"

	// required for database requests.
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/storage/entity"
)

type SQLStorage struct {
	connPool sqlx.DB
}

func New(config config.DatabaseConf) SQLStorage {
	storage := SQLStorage{}
	storage.Connect(config)

	return storage
}

func (s *SQLStorage) Connect(config config.DatabaseConf) error {
	db, err := sqlx.Connect("mysql", "root:123456@(localhost:3306)/calendar?parseTime=true")
	if err != nil {
		return err
	}

	s.connPool = *db

	return nil
}

func (s *SQLStorage) Close() error {
	return s.connPool.Close()
}

func (s SQLStorage) CreateEvent(event entity.Event) error {
	sql := `INSERT INTO events (user_id, title, event_time, duration_min, notice_before, description)
 		VALUES (:user_id, :title, :event_time, :duration_min, :notice_before, :description)`

	data := map[string]interface{}{
		"user_id":       event.UserID,
		"title":         event.Title,
		"event_time":    event.EventTime,
		"duration_min":  event.DurationMin,
		"notice_before": event.NoticeBefore,
		"description":   event.Description,
	}

	_, err := s.connPool.NamedExec(sql, data)
	defer s.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s SQLStorage) UpdateEvent(eventID int, event entity.Event) error {
	v := reflect.ValueOf(event)

	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}

	sql := `UPDATE events SET user_id = :user_id, title = :title, event_time = :event_time,
                  duration_min := duration_min, notice_before = :notice_before, description = :description
                  WHERE id = :event_id`

	data := map[string]interface{}{
		"event_id":      eventID,
		"user_id":       event.UserID,
		"title":         event.Title,
		"event_time":    event.EventTime,
		"duration_min":  event.DurationMin,
		"notice_before": event.NoticeBefore,
		"description":   event.Description,
	}

	defer s.Close()

	_, err := s.connPool.NamedExec(sql, data)
	if err != nil {
		return err
	}

	return nil
}

func (s SQLStorage) RemoveEvent(eventID int) error {
	sql := `DELETE FROM events WHERE id = ?`

	_, err := s.connPool.Exec(sql, eventID)
	defer s.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s SQLStorage) GetEventsList(userID int) ([]entity.Event, error) {
	sql := `SELECT * FROM events WHERE user_id = :userId`

	rows, err := s.connPool.NamedQuery(sql, map[string]interface{}{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Event
	event := entity.Event{}

	for rows.Next() {
		err := rows.StructScan(&event)
		if err != nil {
			return nil, err
		}

		result = append(result, event)
	}

	return result, nil
}
