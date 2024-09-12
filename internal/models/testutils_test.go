package models

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func newTestDB(t *testing.T) *sql.DB {
	// Открываем временную базу данных SQLite в памяти.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	// Загружаем и выполняем setup.sql скрипт.
	setupSQL, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(setupSQL))
	if err != nil {
		t.Fatal(err)
	}

	// Регистрируем функцию очистки после теста.
	t.Cleanup(func() {
		teardownSQL, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(teardownSQL))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	})

	return db
}
