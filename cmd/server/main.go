package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"

	"mshocrylas.com/music-league/internal/db"
	"mshocrylas.com/music-league/internal/schema"
)

type app struct {
	q *db.Queries
}

func main() {
	conn, err := sql.Open("sqlite", "file:musicleague.db?mode=rwc&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	conn.SetMaxOpenConns(1) // avoids "database is locked" ??

	if _, err := conn.ExecContext(context.Background(), schema.SQL); err != nil {
		log.Fatal(err)
	}

	a := &app{q: db.New(conn)}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /rounds/{id}", a.showRound)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("listening on:", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func (a *app) showRound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("round " + r.PathValue("id")))
}
