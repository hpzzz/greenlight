package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Nest json response in another field?
type envelope map[string]interface{}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	return id, err

}

// Change the data parameter to have the type envelope instead of interface
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// encode data to JSON, returning the error if there was one
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a new line to make it easier to view in terminal applciations
	js = append(js, '\n')

	// Loop through the header map and add each header to the http.responsewriter header map
	// note that it's ok if the provided header map is nil/ Go doesnt throw an error if your try to range over a nil map
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil

}
