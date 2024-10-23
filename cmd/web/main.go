package main

import (
	"flag"
	"log/slog"
	"os"

	"forum/handler"
	"forum/repository"
	"forum/service"

	// Go вызывает функцию init() внутри этого пакета.
	_ "forum/internal/memory"
	"forum/internal/session"

	_ "github.com/mattn/go-sqlite3"
)

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

	db, err := repository.NewSqliteDB(*dsn)
	if err != nil {
		logger.Error("Unable to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Создание таблиц и добавление тестовых данных
	if err := repository.InitSqliteDB(db); err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	templateCache, err := handler.NewTemplateCache()
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

	app := &handler.Application{
		Logger:         logger,
		Service:        service.NewService(repository.NewRepository(db)),
		TemplateCache:  templateCache,
		SessionManager: sessionManager,
	}

	err = app.Serve(addr)
	if err != nil {
		logger.Error("Fatal server error", "error", err)
	}
}
