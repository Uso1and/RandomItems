package repositories

import (
	"RandomItems/internal/domain/models"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateDropEvent(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewDropRepository(db)

	now := time.Now()

	dropEvent := &models.DropEvent{
		UserID:       1,
		ItemID:       2,
		DroppedAt:    now,
		IsGuaranteed: true,
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO drop_events").
			WithArgs(dropEvent.UserID, dropEvent.ItemID, dropEvent.DroppedAt, dropEvent.IsGuaranteed).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.CreateDropEvent(context.Background(), dropEvent)
		assert.NoError(t, err)
		assert.Equal(t, 1, dropEvent.ID)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO drop_events").
			WithArgs(dropEvent.UserID, dropEvent.ItemID, dropEvent.DroppedAt, dropEvent.IsGuaranteed).
			WillReturnError(sql.ErrConnDone)

		err := repo.CreateDropEvent(context.Background(), dropEvent)

		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)

	})

	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestGetUserDropHistory(t *testing.T) {

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewDropRepository(db)

	now := time.Now()

	expectedEvent := []*models.DropEvent{
		{
			ID:           1,
			UserID:       1,
			ItemID:       2,
			DroppedAt:    now,
			IsGuaranteed: true,
		},
		{
			ID:           2,
			UserID:       1,
			ItemID:       3,
			DroppedAt:    now.Add(-time.Hour),
			IsGuaranteed: false,
		},
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "item_id", "dropped_at", "is_guaranteed"}).
			AddRow(expectedEvent[0].ID, expectedEvent[0].UserID, expectedEvent[0].ItemID, expectedEvent[0].DroppedAt,
				expectedEvent[0].IsGuaranteed).
			AddRow(expectedEvent[1].ID, expectedEvent[1].UserID, expectedEvent[1].ItemID, expectedEvent[1].DroppedAt,
				expectedEvent[1].IsGuaranteed)

		mock.ExpectQuery("SELECT id, user_id, item_id, dropped_at, is_guaranteed FROM drop_events").
			WithArgs(1, 2).
			WillReturnRows(rows)

		events, err := repo.GetUserDropHistory(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.Equal(t, expectedEvent, events)
	})

	t.Run("error", func(t *testing.T) {

		mock.ExpectQuery("SELECT id, user_id, item_id, dropped_at, is_guaranteed FROM drop_events").
			WithArgs(1, 2).
			WillReturnError(sql.ErrNoRows)

		events, err := repo.GetUserDropHistory(context.Background(), 1, 2)

		assert.Error(t, err)

		assert.Nil(t, events)
	})

	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestGetLastDropTime(t *testing.T) {

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewDropRepository(db)

	now := time.Now()

	t.Run("success", func(t *testing.T) {

		mock.ExpectQuery("SELECT dropped_at FROM drop_events").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"dropped_at"}).AddRow(now))

		result, err := repo.GetLastUserDropTime(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, now, result)
	})

	t.Run("success no rows", func(t *testing.T) {
		mock.ExpectQuery("SELECT dropped_at FROM drop_events").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetLastUserDropTime(context.Background(), 1)
		assert.NoError(t, err)
		assert.True(t, result.IsZero())
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT dropped_at FROM drop_events").
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		result, err := repo.GetLastUserDropTime(context.Background(), 1)
		assert.Error(t, err)
		assert.True(t, result.IsZero())
	})
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUserPityCounter(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewDropRepository(db)

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET pity_counter").
			WithArgs(5, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUserPityCounter(context.Background(), 1, 5)

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET pity_counter").
			WithArgs(5, 1).
			WillReturnError(sql.ErrConnDone)

		err := repo.UpdateUserPityCounter(context.Background(), 1, 5)

		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
	})
	assert.NoError(t, mock.ExpectationsWereMet())
}
