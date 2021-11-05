package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.karolharasim.net/internal/data"
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

	// Create a new instance of the Movie struct, containing the ID we extracted from the URL and some dummy data. Also notice we deliberately havent set a value for the Year field
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Lost",
		Runtime:   102,
		Genres:    []string{"Thriller", "drama"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}
