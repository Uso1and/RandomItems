package repositories

import (
	"RandomItems/internal/domain/models"
	"context"
	"database/sql"
)

type UserRepositoryInterface interface {
	CreateUserRep(c context.Context, user *models.User) error
	GetUser(c context.Context, userID int) (*models.User, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUserRep(c context.Context, user *models.User) error {

	query := `INSERT INTO users (username, created_at, pity_counter) VALUES ($1, $2, $3) RETURNING id`

	return r.db.QueryRowContext(
		c,
		query,
		user.Username,
		user.CreatedAt,
		user.PityCounter).Scan(&user.ID)
}

func (r *UserRepository) GetUser(c context.Context, userID int) (*models.User, error) {

	user := &models.User{ID: userID}

	query := `SELECT username, created_at, pity_counter FROM users WHERE id = $1`

	err := r.db.QueryRowContext(c, query, userID).Scan(&user.Username, &user.CreatedAt, &user.PityCounter)

	if err != nil {
		return nil, err
	}

	return user, nil
}
