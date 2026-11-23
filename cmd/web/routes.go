package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *app) routes() http.Handler {
	mux := http.NewServeMux()

	dynamic := alice.New(http.NewCrossOriginProtection().Handler)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home)) // {$} restricts this route to exact matches
	mux.Handle("GET /login", dynamic.ThenFunc(app.userLogin))

	// TODO: add auth
	// protected := dynamic.Append(<authentication>)

	mux.Handle("GET /leagues", dynamic.ThenFunc(app.listLeagues))
	mux.Handle("GET /leagues/{id}", dynamic.ThenFunc(app.showLeague))
	mux.Handle("POST /leagues", dynamic.ThenFunc(app.createLeague))

	mux.Handle("POST /leagues/{id}/rounds", dynamic.ThenFunc(app.createRound))
	mux.Handle("GET /leagues/{id}/rounds", dynamic.ThenFunc(app.listRoundsInLeague))
	mux.Handle("GET /rounds/{id}", dynamic.ThenFunc(app.showRound))

	return mux
}
