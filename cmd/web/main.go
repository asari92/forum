package main

import (
	"database/sql"
	"flag"
	"forum/internal/models"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	posts *models.PostModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "./forum.db", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Создание таблиц и добавление тестовых данных
	if err := initDB(db); err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		posts: &models.PostModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Call the new app.routes() method to get the servemux containing our routes.
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func initDB(db *sql.DB) error {
	query, err := os.ReadFile("./docs/new_forum.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(query))
	if err != nil {
		return err
	}

	// Проверка наличия данных
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM posts")
	err = row.Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Данные уже есть, ничего не делаем
	}

	// Добавление тестовых данных
	insertTestdataSQL, err := os.ReadFile("./docs/testdata.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(insertTestdataSQL))
	if err != nil {
		return err
	}

	return nil
}
