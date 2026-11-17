package main

import "net/http"

func (a *app) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}

func (a *app) userLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("you should log in but you can't yet"))
}

// The round id is globally unique; but we need to check that the 
// current user has the correct permissions to access the round
func (a *app) showRound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("round " + r.PathValue("id")))
}

func (a *app) showLeague(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("league " + r.PathValue("id")))
}
