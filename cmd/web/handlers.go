package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mshocrylas.com/music-league/internal/db"
	"mshocrylas.com/music-league/internal/validator"
	"mshocrylas.com/music-league/internal/views"
)

func (a *app) home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a *app) showLogin(w http.ResponseWriter, r *http.Request) {
	views.Login().Render(r.Context(), w)
}

func (a *app) login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		a.serverError(w, r, err)
		return
	}
	username := r.PostForm.Get("username")
	fmt.Fprintf(w, "Hello %s", username)
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
		a.clientError(w, http.StatusNotFound)
		return
	}

	views.LeagueList(leagues).Render(r.Context(), w)
}

func (a *app) showLeague(w http.ResponseWriter, r *http.Request) {
	league, ok := a.leagueFromRequest(w, r)
	if !ok {
		return
	}

	a.renderLeague(w, r, http.StatusOK, league, views.RoundForm{})
}

func (a *app) renderLeague(w http.ResponseWriter, r *http.Request, status int, league db.League, form views.RoundForm) {
	rounds, err := a.q.ListRoundsByLeague(r.Context(), league.LeagueID)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	// Call write header just before rendering, so it is after any possible errors that would want to set the header instead.
	w.WriteHeader(status)
	views.League(views.LeaguePage{
		League: league,
		Rounds: rounds,
		Form:   form,
	}).Render(r.Context(), w)
}

func (a *app) createRound(w http.ResponseWriter, r *http.Request) {
	league, ok := a.leagueFromRequest(w, r)
	if !ok {
		return
	}

	if err := r.ParseForm(); err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	// TODO: derive other fields from the user form then you dont need this rounds query even
	rounds, err := a.q.ListRoundsByLeague(r.Context(), league.LeagueID)
	if err != nil {
		a.serverError(w, r, err)
		return
	}
	form := views.RoundForm{
		Theme:              strings.TrimSpace(r.PostForm.Get("theme")),
		Position:           int64(len(rounds) + 1),
		SubmissionDeadline: time.Now(),
		VotingDeadline:     time.Now(),
	}

	// TODO: add checks on the other values
	form.CheckField(validator.NotBlank(form.Theme), "theme", "Theme is required")

	if !form.Valid() {
		a.renderLeague(w, r, http.StatusUnprocessableEntity, league, form)
		return
	}

	if _, err := a.q.CreateRound(r.Context(), db.CreateRoundParams{
		LeagueID:           league.LeagueID,
		Theme:              form.Theme,
		Position:           form.Position,
		SubmissionDeadline: dateToISOString(form.SubmissionDeadline),
		VotingDeadline:     dateToISOString(form.VotingDeadline),
	}); err != nil {
		a.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/leagues/%d", league.LeagueID), http.StatusSeeOther)
}

// The round id is globally unique; but we need to check that the
// current user has the correct permissions to access the round
// * This page should show different things depending on the "phase" of the round.
// Submitting, Voting, or Results
// "Submitting" shows nothing except a form for the user to submit a song for the round
// once the user has submitted, this is indicated. Maybe they can revise their choice while
// the round is still waiting on others to submit.
// "Voting" shows everyone's submitted song choice, ideally in some random/anonymized order
// The page should be a form where the user ranks every choice from favorite to least.
// Once they submit their votes, they can probably see a results-in-progress sort of page.
// Or maybe not, TBD.
// TODO: Add the ability to leave a comment with a vote, that's like the best part.
// "Results" after everyone has voted, final results are tallied and the results page is shown.
// All submissions are displayed ranked from highest to lowest rank, along with how each person
// voted on the individual songs and what their comment was if they left one.
func (a *app) showRound(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		a.clientError(w, http.StatusUnprocessableEntity)
		return
	}
	round, err := a.q.GetRound(r.Context(), id)
	if err != nil {
		a.clientError(w, http.StatusNotFound)
		return
	}

	views.Round(round).Render(r.Context(), w)
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

func (a *app) newSong(w http.ResponseWriter, r *http.Request) {
	views.NewSong(views.SongForm{}).Render(r.Context(), w)
}

func (a *app) createSong(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form := views.SongForm{
		Title:  strings.TrimSpace(r.PostForm.Get("title")),
		Artist: strings.TrimSpace(r.PostForm.Get("artist")),
		URL:    strings.TrimSpace(r.PostForm.Get("url")),
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "Title is required")
	form.CheckField(validator.NotBlank(form.Artist), "artist", "Artist is required")
	if validator.NotBlank(form.URL) {
		form.CheckField(validator.ValidHTTPURL(form.URL), "url", "Enter a link starting with http:// or https://")
	}

	if !form.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		views.NewSong(form).Render(r.Context(), w)
		return
	}

	link := sql.NullString{String: form.URL, Valid: form.URL != ""}
	song, err := a.q.CreateSong(r.Context(), db.CreateSongParams{Title: form.Title, Artist: form.Artist, Url: link})
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	redirectURL := fmt.Sprintf("/songs/%d", song.SongID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (a *app) showSong(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		a.clientError(w, http.StatusUnprocessableEntity)
		return
	}
	song, err := a.q.GetSong(r.Context(), id)
	if err != nil {
		a.clientError(w, http.StatusNotFound)
		return
	}

	views.Song(song).Render(r.Context(), w)
}

func (a *app) leagueFromRequest(w http.ResponseWriter, r *http.Request) (db.League, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		a.clientError(w, http.StatusNotFound)
		return db.League{}, false
	}
	league, err := a.q.GetLeague(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.clientError(w, http.StatusNotFound)
		} else {
			a.serverError(w, r, err)
		}
		return db.League{}, false
	}
	return league, true
}
