package repositories

import (
	"RandomItems/internal/domain/models"
	"context"
	"database/sql"
)

type ItemRepositoryInterface interface {
	GetItems(c context.Context) ([]*models.Item, error)
	GetItem(c context.Context, id int) (*models.Item, error)
}

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) GetItems(c context.Context) ([]*models.Item, error) {
	query := `SELECT id, name, rarity, base_chance, min_pity FROM items`
	rows, err := r.db.QueryContext(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Rarity, &item.BaseChance, &item.MinPity); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (r *ItemRepository) GetItem(c context.Context, id int) (*models.Item, error) {

	query := `SELECT id, name, rarity, base_chance, min_pity FROM items WHERE id = $1`

	item := &models.Item{ID: id}

	err := r.db.QueryRowContext(c, query, id).Scan(&item.ID, &item.Name, &item.Rarity, &item.BaseChance, &item.MinPity)

	if err != nil {
		return nil, err
	}
	return item, nil
}
