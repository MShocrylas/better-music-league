package main

import "net/http"

func (a *app) showRound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("round " + r.PathValue("id")))
}
