package repositories

import (
	"RandomItems/internal/domain/models"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetItems(t *testing.T) {

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewItemRepository(db)

	expectItems := []*models.Item{

		{
			ID:         1,
			Name:       "Health Potion",
			Rarity:     "common",
			BaseChance: 0.5,
			MinPity:    0,
		},

		{
			ID:         2,
			Name:       "Dragon Sword",
			Rarity:     "legendary",
			BaseChance: 0.05,
			MinPity:    50,
		},
	}

	t.Run("success - get all items", func(t *testing.T) {

		rows := sqlmock.NewRows([]string{"id", "name", "rarity", "base_chance", "min_pity"}).
			AddRow(expectItems[0].ID, expectItems[0].Name, expectItems[0].Rarity, expectItems[0].BaseChance,
				expectItems[0].MinPity).
			AddRow(expectItems[1].ID, expectItems[1].Name, expectItems[1].Rarity, expectItems[1].BaseChance,
				expectItems[1].MinPity)

		mock.ExpectQuery("SELECT id, name, rarity, base_chance, min_pity FROM items").
			WillReturnRows(rows)

		items, err := repo.GetItems(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectItems, items)
	})

	t.Run("error - database error", func(t *testing.T) {

		mock.ExpectQuery("SELECT id, name, rarity, base_chance, min_pity").
			WillReturnError(sql.ErrConnDone)

		items, err := repo.GetItems(context.Background())

		assert.Error(t, err)
		assert.Nil(t, items)
		assert.Equal(t, sql.ErrConnDone, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetItem(t *testing.T) {

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewItemRepository(db)

	expectedItem := &models.Item{

		ID:         1,
		Name:       "Health Potion",
		Rarity:     "common",
		BaseChance: 0.5,
		MinPity:    0,
	}

	t.Run("success - get user by id ", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "rarity", "base_chance", "min_pity"}).
			AddRow(expectedItem.ID, expectedItem.Name, expectedItem.Rarity, expectedItem.BaseChance, expectedItem.MinPity)

		mock.ExpectQuery("SELECT id, name, rarity, base_chance, min_pity FROM items WHERE id = ?").
			WithArgs(1).
			WillReturnRows(rows)

		item, err := repo.GetItem(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedItem, item)
	})

	t.Run("error - item not found", func(t *testing.T) {

		mock.ExpectQuery("SELECT id, name, rarity, base_chance, min_pity FROM items WHERE id = ?").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		item, err := repo.GetItem(context.Background(), 999)

		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("error - database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, rarity, base_chance, min_pity FROM items WHERE id = ?").
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		item, err := repo.GetItem(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Equal(t, sql.ErrConnDone, err)
	})
	assert.NoError(t, mock.ExpectationsWereMet())
}
