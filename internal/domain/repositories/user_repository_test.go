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

func TestCreateUser(t *testing.T) {

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("Ошибка создания mock БД: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)
	now := time.Now()

	user := &models.User{
		Username:    "testuser",
		CreatedAt:   now,
		PityCounter: 0,
	}

	t.Run("success - user creation", func(t *testing.T) {

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Username, user.CreatedAt, user.PityCounter).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.CreateUserRep(context.Background(), user)

		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
	})

	t.Run("error - when creating a user", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Username, user.CreatedAt, user.PityCounter).
			WillReturnError(sql.ErrConnDone)

		err := repo.CreateUserRep(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
	})
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUser(t *testing.T) {

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("Ошибка создания mock БД: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	now := time.Now()

	expectedUser := &models.User{

		ID:          1,
		Username:    "testuser",
		CreatedAt:   now,
		PityCounter: 5,
	}

	t.Run("success - get user", func(t *testing.T) {

		rows := sqlmock.NewRows([]string{"username", "created_at", "pity_counter"}).
			AddRow(expectedUser.Username, expectedUser.CreatedAt, expectedUser.PityCounter)

		mock.ExpectQuery("SELECT username, created_at, pity_counter FROM users WHERE id = ?").
			WithArgs(1).
			WillReturnRows(rows)

		user, err := repo.GetUser(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("error - user not found", func(t *testing.T) {

		mock.ExpectQuery("SELECT username, created_at, pity_counter FROM users WHERE id = ?").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUser(context.Background(), 999)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("error - database", func(t *testing.T) {
		mock.ExpectQuery("SELECT username, created_at, pity_counter FROM users WHERE id = ?").
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		user, err := repo.GetUser(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, sql.ErrConnDone, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
