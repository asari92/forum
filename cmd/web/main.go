package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"forum/docs"

	// Go вызывает функцию init() внутри этого пакета.
	_ "forum/internal/memory"
	"forum/internal/models"
	"forum/internal/session"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	users          models.UserModelInterface
	posts          models.PostModelInterface
	categories     models.CategoryModelInterface
	templateCache  map[string]*template.Template
	sessionManager *session.Manager
}

func main() {
	addr := flag.String("addr", ":7777", "HTTP network address")
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

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager, err := session.NewManager("memory", "gosessionid", 600)
	if err != nil {
		errorLog.Fatal(err)
	}

	go sessionManager.GC()

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		users:          &models.UserModel{DB: db},
		posts:          &models.PostModel{DB: db},
		categories:     &models.CategoryModel{DB: db},
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}

	// Чтение встроенных TLS-ключей из файловой системы

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on https://localhost%s", *addr)
	err = srv.ListenAndServe() // Пустые строки означают, что сертификат и ключ уже загружены в `tlsConfig`
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
	// Используем встроенный файл для инициализации базы данных
	query, err := docs.Files.ReadFile("new_forum.sql")
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
	insertTestdataSQL, err := docs.Files.ReadFile("testdata.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(insertTestdataSQL))
	if err != nil {
		return err
	}

	return nil
}
