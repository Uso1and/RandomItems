package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	// Создаем временную директорию для миграций
	tempDir := t.TempDir()
	migrationsDir := filepath.Join(tempDir, "migrations")
	err := os.Mkdir(migrationsDir, 0755)
	require.NoError(t, err)

	// Создаем тестовые файлы миграций
	migrations := []struct {
		name    string
		content string
	}{
		{
			name:    "0001_test.sql",
			content: "CREATE TABLE IF NOT EXISTS test_table (id SERIAL PRIMARY KEY);",
		},
		{
			name:    "0002_test.sql",
			content: "ALTER TABLE test_table ADD COLUMN IF NOT EXISTS name VARCHAR(100);",
		},
	}

	for _, migration := range migrations {
		err = os.WriteFile(filepath.Join(migrationsDir, migration.name), []byte(migration.content), 0644)
		require.NoError(t, err)
	}

	SetMigrationsDir(migrationsDir)

	defer SetMigrationsDir("")

	t.Run("successful initialization", func(t *testing.T) {

		originalConnStr := connStr
		connStr = "user=postgres password=root dbname=test_db sslmode=disable"
		defer func() { connStr = originalConnStr }()

		err := Init()
		assert.NoError(t, err)
		assert.NotNil(t, DB)

		var exists bool
		err = DB.QueryRow(`
            SELECT EXISTS (
                SELECT FROM information_schema.tables 
                WHERE table_name = 'test_table'
            )`).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists)

		_, err = DB.Exec("INSERT INTO test_table (name) VALUES ('test')")
		assert.NoError(t, err)
	})

	t.Run("invalid migration", func(t *testing.T) {
		err := os.WriteFile(
			filepath.Join(migrationsDir, "0003_broken.sql"),
			[]byte("INVALID SQL STATEMENT;"),
			0644,
		)
		require.NoError(t, err)
		err = Init()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "migration 0003_broken.sql failed")
	})
}

func TestApplyMigrations(t *testing.T) {

	db, err := sql.Open("postgres", "user=postgres password=root dbname=test_db sslmode=disable")
	require.NoError(t, err)
	defer db.Close()

	oldDB := DB
	oldMigrationsDir := migrationsDir
	defer func() {
		DB = oldDB
		SetMigrationsDir(oldMigrationsDir)

		db.Exec("DROP TABLE IF EXISTS test_migrations")
	}()

	DB = db

	t.Run("successful migrations", func(t *testing.T) {
		tempDir := t.TempDir()
		migrationsDir := filepath.Join(tempDir, "migrations")
		require.NoError(t, os.Mkdir(migrationsDir, 0755))
		SetMigrationsDir(migrationsDir)

		migrations := []struct {
			name    string
			content string
		}{
			{
				name:    "0001_create_table.sql",
				content: "CREATE TABLE IF NOT EXISTS test_migrations (id SERIAL PRIMARY KEY);",
			},
			{
				name:    "0002_add_column.sql",
				content: "ALTER TABLE test_migrations ADD COLUMN IF NOT EXISTS sample_column TEXT;",
			},
		}

		for _, m := range migrations {
			err := os.WriteFile(filepath.Join(migrationsDir, m.name), []byte(m.content), 0644)
			require.NoError(t, err)
		}

		err := applyMigrations()
		assert.NoError(t, err)

		var exists bool
		err = db.QueryRow(`
            SELECT EXISTS (
                SELECT FROM information_schema.tables 
                WHERE table_name = 'test_migrations'
            )`).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists)

		_, err = db.Exec("INSERT INTO test_migrations (sample_column) VALUES ('test')")
		assert.NoError(t, err)
	})

	t.Run("invalid migration", func(t *testing.T) {
		tempDir := t.TempDir()
		migrationsDir := filepath.Join(tempDir, "migrations")
		require.NoError(t, os.Mkdir(migrationsDir, 0755))
		SetMigrationsDir(migrationsDir)

		err := os.WriteFile(
			filepath.Join(migrationsDir, "0001_broken.sql"),
			[]byte("INVALID SQL COMMAND;"),
			0644,
		)
		require.NoError(t, err)

		err = applyMigrations()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "migration 0001_broken.sql failed")
	})

	t.Run("empty directory", func(t *testing.T) {
		tempDir := t.TempDir()
		migrationsDir := filepath.Join(tempDir, "empty_migrations")
		require.NoError(t, os.Mkdir(migrationsDir, 0755))
		SetMigrationsDir(migrationsDir)

		err := applyMigrations()
		assert.NoError(t, err)
	})
}
