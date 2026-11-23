package main

import (
	"context"
	"database/sql"
	_ "embed"
	"flag"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pressly/goose/v3"

	_ "modernc.org/sqlite"

	"mshocrylas.com/music-league/internal/db"
	"mshocrylas.com/music-league/migrations"
)

type app struct {
	q      *db.Queries
	logger *slog.Logger
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}))

	if err := run(logger); err != nil {
		logger.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	path := flag.String("db", "app.db", "SQLite database path")
	flag.Parse()

	q := url.Values{}
	q.Add("_foreign_keys", "true")
	q.Add("_busy_timeout", "5000")
	q.Add("_journal_mode", "WAL")
	q.Add("_synchronous", "NORMAL")
	q.Add("_txlock", "IMMEDIATE")

	dsn := "file:" + *path + "?" + q.Encode()
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Test the connection before using it for anything
	if err := conn.Ping(); err != nil {
		return err
	}

	if err := migrate(conn, logger); err != nil {
		return err
	}

	a := &app{q: db.New(conn), logger: logger}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      a.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.Info("Server listening", "addr", srv.Addr)
	return srv.ListenAndServe()
}

func migrate(db *sql.DB, logger *slog.Logger) error {
	gooseLogger := &gooseSlogAdapter{
		logger: logger.With(slog.String("component", "goose")),
	}

	provider, err := goose.NewProvider(
		goose.DialectSQLite3,
		db,
		migrations.FS,
		goose.WithLogger(gooseLogger),
		goose.WithVerbose(true),
	)
	if err != nil {
		return err
	}
	_, err = provider.Up(context.Background())
	return err
}
