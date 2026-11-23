package main

import (
	"net/http"
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
