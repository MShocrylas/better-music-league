package main

import (
	"net/http"
	"strconv"
)

func (a *app) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}

func (a *app) userLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("you should log in but you can't yet"))
}

func (a *app) createLeague(w http.ResponseWriter, r *http.Request) {
	leagueID, err := a.q.CreateLeague(r.Context(), "test")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/leagues/"+strconv.Itoa(int(leagueID.LeagueID)), http.StatusSeeOther)
}

func (a *app) listLeagues(w http.ResponseWriter, r *http.Request) {
	leagues, err := a.q.ListLeagues(r.Context())
	if err != nil {
		a.logger.Error("oops", "error", err)
		http.NotFound(w, r)
	}

	w.Write([]byte("leagues: "))
	for l := range leagues {
		w.Write([]byte(strconv.Itoa(l) + "\n"))
	}
}

func (a *app) showLeague(w http.ResponseWriter, r *http.Request) {
	intID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	league, err := a.q.GetLeague(r.Context(), intID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	w.Write([]byte("league: "))
	w.Write([]byte(strconv.Itoa(int(league.LeagueID)) + " " + league.Name))
}

func (a *app) listRounds(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("round " + r.PathValue("id")))
}

// The round id is globally unique; but we need to check that the
// current user has the correct permissions to access the round
func (a *app) showRound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("round " + r.PathValue("id")))
}
