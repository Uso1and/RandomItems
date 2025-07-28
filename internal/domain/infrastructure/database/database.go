package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"sort"

	_ "github.com/lib/pq"
)

var (
	DB            *sql.DB
	connStr       = "user=postgres password=root dbname=ridb sslmode=disable"
	migrationsDir = ""
)

// Для тестов
func SetMigrationsDir(dir string) {
	migrationsDir = dir
}

func Init() error {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	if err = applyMigrations(); err != nil {
		return fmt.Errorf("migration failed: %v", err)
	}

	return nil
}

func applyMigrations() error {
	dir := migrationsDir
	if dir == "" {
		_, filename, _, _ := runtime.Caller(0)
		dir = filepath.Join(filepath.Dir(filename), "..", "migrations")
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return err
		}

		if _, err := DB.Exec(string(content)); err != nil {
			return fmt.Errorf("migration %s failed: %v", file.Name(), err)
		}
	}

	return nil
}
