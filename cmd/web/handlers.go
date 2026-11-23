package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"mshocrylas.com/music-league/internal/db"
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
		a.serverError(w, r, err)
		return
	}

	redirectURL := fmt.Sprintf("/leagues/%d", leagueID.LeagueID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (a *app) listLeagues(w http.ResponseWriter, r *http.Request) {
	leagues, err := a.q.ListLeagues(r.Context())
	if err != nil {
		a.logger.Error("oops", "error", err)
		a.clientError(w, http.StatusNotFound)
		return
	}

	w.Write([]byte("leagues: "))
	for _, l := range leagues {
		w.Write([]byte(strconv.Itoa(int(l.LeagueID)) + "\n"))
	}
}

func (a *app) showLeague(w http.ResponseWriter, r *http.Request) {
	intID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		a.serverError(w, r, err)
		return
	}
	league, err := a.q.GetLeague(r.Context(), intID)
	if err != nil {
		a.clientError(w, http.StatusNotFound)
		return
	}
	w.Write([]byte("league: "))
	w.Write([]byte(strconv.Itoa(int(league.LeagueID)) + " " + league.Name))
}

func (a *app) createRound(w http.ResponseWriter, r *http.Request) {
	leagueID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		a.clientError(w, http.StatusUnprocessableEntity)
		return
	}
	// TODO: VERY TEMPORARY
	round, err := a.q.CreateRound(r.Context(), db.CreateRoundParams{
		LeagueID:           leagueID,
		Theme:              "theme",
		Position:           1,
		SubmissionDeadline: sql.NullString{},
		VotingDeadline:     sql.NullString{},
	})
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	fmt.Fprintf(w, "Round created: %d %s\n", round.RoundID, round.Theme)
}

func (a *app) listRoundsInLeague(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		a.clientError(w, http.StatusUnprocessableEntity)
		return
	}
	rounds, err := a.q.ListRoundsByLeague(r.Context(), id)
	if err != nil {
		a.serverError(w, r, err)
		return
	}
	for _, round := range rounds {
		fmt.Fprintf(w, "Round %d: %s\n", round.RoundID, round.Theme)
	}
}

// The round id is globally unique; but we need to check that the
// current user has the correct permissions to access the round
func (a *app) showRound(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		a.clientError(w, http.StatusUnprocessableEntity)
		return
	}
	round, err := a.q.GetRound(r.Context(), id)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	fmt.Fprintf(w, "Round %d: %s in league %d\n", round.RoundID, round.Theme, round.LeagueID)
}
