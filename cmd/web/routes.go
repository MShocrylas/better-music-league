package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *app) routes() http.Handler {
	mux := http.NewServeMux()

	dynamic := alice.New(http.NewCrossOriginProtection().Handler)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home)) // {$} restricts this route to exact matches
	mux.Handle("GET /login", dynamic.ThenFunc(app.showLogin))
	mux.Handle("POST /login", dynamic.ThenFunc(app.login))

	// TODO: add auth
	// protected := dynamic.Append(<authentication>)

	mux.Handle("GET /leagues", dynamic.ThenFunc(app.listLeagues))
	mux.Handle("POST /leagues", dynamic.ThenFunc(app.createLeague))
	mux.Handle("GET /leagues/{id}", dynamic.ThenFunc(app.showLeague))

	mux.Handle("GET /leagues/{id}/rounds", dynamic.ThenFunc(app.listRoundsInLeague))
	mux.Handle("POST /leagues/{id}/rounds", dynamic.ThenFunc(app.createRound))
	mux.Handle("GET /rounds/{id}", dynamic.ThenFunc(app.showRound))

	mux.Handle("POST /songs", dynamic.ThenFunc(app.createSong))
	mux.Handle("GET /songs/new", dynamic.ThenFunc(app.newSong))
	mux.Handle("GET /songs/{id}", dynamic.ThenFunc(app.showSong))

	return mux
}
