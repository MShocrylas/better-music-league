package main

import (
	"context"
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/pressly/goose/v3"

	_ "modernc.org/sqlite"

	"mshocrylas.com/music-league/internal/db"
	"mshocrylas.com/music-league/migrations"
)

type gooseSlogAdapter struct {
	logger *slog.Logger
}

func (s *gooseSlogAdapter) Printf(format string, v ...any) {
	msg := strings.TrimRight(fmt.Sprintf(format, v...), "\n")
	s.logger.Info(msg)
}

func (s *gooseSlogAdapter) Fatalf(format string, v ...any) {
	msg := strings.TrimRight(fmt.Sprintf(format, v...), "\n")
	s.logger.Error(msg)
	os.Exit(1)
}

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
	q.Add("_pragma", "foreign_keys(1)")
	q.Add("_pragma", "busy_timeout(5000)")
	q.Add("_pragma", "journal_mode(WAL)")
	q.Add("_pragma", "synchronous(NORMAL)")

	dsn := "file:" + *path + "?" + q.Encode()
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.SetMaxOpenConns(1) // avoids "database is locked" ??

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
	logger.Info("listening on:", "addr", srv.Addr)
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
