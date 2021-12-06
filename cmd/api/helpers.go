package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError

	switch {
	// Use the errors.As() function to check whether the error has the type
	// *json.SyntaxError. If it does, then return a plain-english error message
	// which includes the location of the problem.”
	case errors.As(err, &syntaxError):
		return fmt.Errorf("body contains badly-formated JSON (at character %d)", syntaxError.Offset)

	case errors.Is(err, io.ErrUnexpectedEOF):
		return errors.New("body contains badly-formated JSON")
	case errors.As(err, &unmarshalTypeError):
		if unmarshalTypeError.Field != "" {
			return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
		}
		return fmt.Errorf("body contains incorrect json type (at character %d)", unmarshalTypeError.Offset)
	case errors.Is(err, io.EOF):
		return errors.New("body must not be empty")

	case errors.As(err, &invalidUnmarshalError):
		panic(err)
	default:
		return err
	}

}
