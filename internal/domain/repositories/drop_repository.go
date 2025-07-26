package repositories

import (
	"RandomItems/internal/domain/models"
	"context"
	"database/sql"
	"time"
)

type DropRepositoryInterface interface {
	CreateDropEvent(c context.Context, event *models.DropEvent) error
	GetUserDropHistory(c context.Context, userID int, limit int) ([]*models.DropEvent, error)
	GetLastUserDropTime(c context.Context, userID int) (time.Time, error)
	UpdateUserPityCounter(c context.Context, userID int, couter int) error
}
type DropRepository struct {
	db *sql.DB
}

func NewDropRepository(db *sql.DB) *DropRepository {
	return &DropRepository{db: db}
}

func (r *DropRepository) CreateDropEvent(c context.Context, event *models.DropEvent) error {
	query := `INSERT INTO drop_events (user_id, item_id, dropped_at, is_guaranteed) VALUES ($1, $2, $3, $4) RETURNING id`

	return r.db.QueryRowContext(c, query, event.UserID, event.ItemID, event.DroppedAt, event.IsGuaranteed).Scan(&event.ID)
}

func (r *DropRepository) GetUserDropHistory(c context.Context, userID int, limit int) ([]*models.DropEvent, error) {
	query := `SELECT id, user_id, item_id, dropped_at, is_guaranteed FROM drop_events WHERE user_id = $1 ORDER BY dropped_at DESC LIMIT $2`

	rows, err := r.db.QueryContext(c, query, userID, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.DropEvent

	for rows.Next() {
		var event models.DropEvent
		if err := rows.Scan(&event.ID, &event.UserID, &event.ItemID, &event.DroppedAt, &event.IsGuaranteed); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

func (r *DropRepository) GetLastUserDropTime(c context.Context, userID int) (time.Time, error) {
	query := `SELECT dropped_at FROM drop_events WHERE user_id = $1 ORDER BY dropped_at DESC LIMIT 1`

	var lastDrop time.Time

	err := r.db.QueryRowContext(c, query, userID).Scan(&lastDrop)

	if err == sql.ErrNoRows {
		return time.Time{}, nil
	}
	if err != nil {
		return time.Time{}, nil
	}
	return lastDrop, nil
}

func (r *DropRepository) UpdateUserPityCounter(c context.Context, userID int, couter int) error {
	query := `UPDATE users SET pity_counter = $1 WHERE id = $2`
	_, err := r.db.ExecContext(c, query, couter, userID)
	return err
}
