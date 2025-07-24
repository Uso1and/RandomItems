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

var DB *sql.DB

func Init() error {

	connStr := `user=postgres password=root dbname=ridb sslmode=disable`

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

	_, filename, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filename)

	migrationsPath := filepath.Join(basePath, "..", "migrations")

	files, err := ioutil.ReadDir(migrationsPath)
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
		content, err := ioutil.ReadFile(filepath.Join(migrationsPath, file.Name()))
		if err != nil {
			return err
		}
		if _, err := DB.Exec(string(content)); err != nil {
			return fmt.Errorf("migration %s failed: %v", file.Name(), err)
		}
		fmt.Printf("Applied migration %s\n", file.Name())
	}
	return nil
}
