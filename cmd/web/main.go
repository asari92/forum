package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	// Go вызывает функцию init() внутри этого пакета.
	_ "forum/internal/memory"
	"forum/internal/models"
	"forum/internal/session"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	users          *models.UserModel
	posts          *models.PostModel
	categories     *models.CategoryModel
	templateCache  map[string]*template.Template
	sessionManager *session.Manager
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

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Call the new app.routes() method to get the servemux containing our routes.
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
