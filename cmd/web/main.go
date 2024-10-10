package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"forum/docs"
	tlsecurity "forum/tls"

	// Go вызывает функцию init() внутри этого пакета.
	_ "forum/internal/memory"
	"forum/internal/models"
	"forum/internal/session"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	logger           *slog.Logger
	users            models.UserModelInterface
	posts            models.PostModelInterface
	comments         models.CommentModelInterface
	postReactions    models.PostReactionModelInterface
	commentReactions models.CommentReactionModelInterface
	categories       models.CategoryModelInterface
	templateCache    map[string]*template.Template
	sessionManager   *session.Manager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "./forum.db", "MySQL data source name")

	flag.Parse()

	// Инициализация нового логгера slog
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,            // Включить вывод источника вызова (файл и строка)
		Level:     slog.LevelDebug, // задан дебаг уровень, можно поменять на инфо чтобы убрать лишнюю инфу
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				// Устанавливаем формат времени на "2006-01-02 15:04:05"
				a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05"))
			}
			return a
		},
	}))

	slog.SetDefault(logger)

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error("Unable to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Создание таблиц и добавление тестовых данных
	if err := initDB(db); err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error("Failed to create template cache", "error", err)
		os.Exit(1)
	}

	sessionManager, err := session.NewManager("memory", "gosessionid", 600)
	if err != nil {
		logger.Error("Failed to create session manager", "error", err)
		os.Exit(1)
	}

	go sessionManager.GC()

	app := &application{
		logger:           logger,
		users:            &models.UserModel{DB: db},
		posts:            &models.PostModel{DB: db},
		comments:         &models.CommentModel{DB: db},
		postReactions:    &models.PostReactionModel{DB: db},
		commentReactions: &models.CommentReactionModel{DB: db},
		categories:       &models.CategoryModel{DB: db},
		templateCache:    templateCache,
		sessionManager:   sessionManager,
	} 

	// Чтение встроенных TLS-ключей из файловой системы
	certPEM, err := fs.ReadFile(tlsecurity.Files, "cert.pem")
	if err != nil {
		app.logger.Error("Failed to read TLS certificate", "error", err)
		os.Exit(1)
	}
	keyPEM, err := fs.ReadFile(tlsecurity.Files, "key.pem")
	if err != nil {
		app.logger.Error("Failed to read TLS key", "error", err)
		os.Exit(1)
	}
	// Загрузка ключей в формате x509
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		app.logger.Error("Failed to load TLS certificate pair", "error", err)
		os.Exit(1)
	}

	tlsConfig := &tls.Config{
		Certificates:     []tls.Certificate{tlsCert},
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
		Addr:         *addr,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting server", "address", "https://localhost"+*addr)
	err = srv.ListenAndServeTLS("", "") // Пустые строки означают, что сертификат и ключ уже загружены в `tlsConfig`
	if err != nil {
		logger.Error("Server failed to start", "error", err)
	}
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
