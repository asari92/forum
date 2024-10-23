package handler

import (
	"context"
	"crypto/tls"
	"errors"
	"forum/internal/session"
	"forum/service"
	tlsecurity "forum/tls"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	Service        *service.Service
	Logger         *slog.Logger
	TemplateCache  map[string]*template.Template
	SessionManager *session.Manager
}

func NewHandler(usecases *service.Service) *Application {
	return &Application{Service: usecases}
}

func (app *Application) Serve(addr *string) error {
	// Чтение встроенных TLS-ключей из файловой системы
	certPEM, err := fs.ReadFile(tlsecurity.Files, "cert.pem")
	if err != nil {
		app.Logger.Error("Failed to read TLS certificate", "error", err)
		os.Exit(1)
	}
	keyPEM, err := fs.ReadFile(tlsecurity.Files, "key.pem")
	if err != nil {
		app.Logger.Error("Failed to read TLS key", "error", err)
		os.Exit(1)
	}
	// Загрузка ключей в формате x509
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		app.Logger.Error("Failed to load TLS certificate pair", "error", err)
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
		ErrorLog:     slog.NewLogLogger(app.Logger.Handler(), slog.LevelError),
		Handler:      app.Routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  9 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Create a shutdownError channel. We will use this to receive any errors returned
	// by the graceful Shutdown() function.
	shutdownError := make(chan error)

	go func() {
		// Intercept the signals, as before.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		// Update the log entry to say "shutting down server" instead of "caught signal".
		app.Logger.Info("shutting down server", "signal", s.String())

		// Create a context with a 20-second timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// Call Shutdown() on our server, passing in the context we just made.
		// Shutdown() will return nil if the graceful shutdown was successful, or an
		// error (which may happen because of a problem closing the listeners, or
		// because the shutdown didn't complete before the 20-second context deadline is
		// hit). We relay this return value to the shutdownError channel.
		shutdownError <- srv.Shutdown(ctx)
	}()

	app.Logger.Info("Starting server", "address", "https://localhost"+*addr)
	// Calling Shutdown() on our server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error. So if we see this error, it is actually a
	// good thing and an indication that the graceful shutdown has started. So we check
	// specifically for this, only returning the error if it is NOT http.ErrServerClosed.
	err = srv.ListenAndServeTLS("", "") // Пустые строки означают, что сертификат и ключ уже загружены в `tlsConfig`
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Otherwise, we wait to receive the return value from Shutdown() on the
	// shutdownError channel. If return value is an error, we know that there was a
	// problem with the graceful shutdown and we return the error.
	err = <-shutdownError
	if err != nil {
		return err
	}

	// At this point we know that the graceful shutdown completed successfully and we
	// log a "stopped server" message.
	app.Logger.Info("stopped server", "addr", *addr)

	return nil
}
