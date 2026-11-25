package main

import (
	"database/sql"
	"net/http"
	"time"
)

// Logs the error, and sends a generic 500 internal server error response
func (a *app) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		//trace = string(debug.Stack())
	)
	a.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Sends the specified status code with the corresponding description
func (a *app) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// dateToISOString converts a Go time.Time into a UTC ISO-8601 (RFC 3339) sql.NullString.
// Zero times are converted to SQL NULL.
// The sql database stores times as this format, so this should be used for conversion.
func dateToISOString(datetime time.Time) sql.NullString {
	if datetime.IsZero() {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: datetime.UTC().Format(time.RFC3339),
		Valid:  true,
	}
}
