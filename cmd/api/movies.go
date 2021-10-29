package main

import (
	"fmt"
	"net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	//When httprouter is parsing a requestm ant interpolated URL parameters will be sytored in the request conext. We can use the ParamsFromContext() fdunction to retrieve slice containing these parameter names and values

	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "show the details of movue %d\n", id)
}
