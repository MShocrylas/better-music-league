package main

import (
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
		Level:     slog.LevelDebug,
		AddSource: true,
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

	if err := migrate(conn); err != nil {
		return err
	}

	a := &app{q: db.New(conn), logger: logger}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /rounds/{id}", a.showRound)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.Info("listening on:", "addr", srv.Addr)
	return srv.ListenAndServe()
}

func (a *app) showRound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("round " + r.PathValue("id")))
}

func migrate(db *sql.DB) error {
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}
	return goose.Up(db, ".")
}
