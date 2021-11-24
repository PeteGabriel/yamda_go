package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"yamda_go/internal/config"
	"yamda_go/internal/data/provider"
	"yamda_go/internal/models"

	"github.com/julienschmidt/httprouter"
)

//TODO generate this automatically at build time
const (
	version = "0.0.1"
	ID      = "id"
)

//Application type contains all dependencies for the top layer of
//the API.
type Application struct {
	log      *log.Logger
	config   *config.Settings
	provider provider.IMovieProvider
}

//ParseId parses the parameter id present in a given
//route params sent via parameters to this function.
//If the parameter id is not found an error will be returned.
//method doesn't use any dependencies from our application struct
//it could just be a regular function, rather than a method on application.
//But in general, I suggest setting up all your application-specific handlers
//and helpers so that they are methods on application.
//It helps maintain consistency in your code structure,
//and also future-proofs the code for when those handlers and helpers change later,
//and they do need access to a dependency.
func (app *Application) ParseId(p httprouter.Params) (int64, error) {
	num := p.ByName(ID)
	id, err := strconv.ParseInt(num, 10, 64)
	if err != nil || id < 1 {
		return -1, errors.New("invalid id parameter from route parameters")
	}
	return id, nil
}

//envelope type. Allow inserting types and self-document them in JSON responses.
type envelope map[string]interface{}

//Helper method for sending JSON responses in case of error. This takes the destination
// http.ResponseWriter, the HTTP status code to send, the data to encode to JSON, and a
// header map containing any additional HTTP headers we want to include in the response.
//By default, the header "Content-Type" is set to "application/problem+json".
func (app *Application) writeError(w http.ResponseWriter, status int, data models.ErrorProblem, headers http.Header) error {
	resp, err := json.Marshal(data)
	if err != nil {
		app.log.Println(err)
		return errors.New("an error happened at the server level")
	}
	//apply all values to respective header keys
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status) //status must be the last write
	w.Write(resp)
	return nil
}

//Helper method for sending JSON responses. This takes the destination
// http.ResponseWriter, the HTTP status code to send, the data to encode to JSON, and a
// header map containing any additional HTTP headers we want to include in the response.
//By default, the header "Content-Type" is set to "application/json".
func (app *Application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	resp, err := json.Marshal(data)
	if err != nil {
		app.log.Println(err)
		return errors.New("an error happened at the server level")
	}
	//apply all values to respective header keys
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status) //status must be the last write
	w.Write(resp)
	return nil
}

//decode body content into JSON content. If any of the validations fails,
//an human redable error as well as an error code is returned by this method.
//Althought the default code is 200, the error code should only be considered in case of error.
func (app *Application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	//limit size of request body to 1MB
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() //do not accept unknown fields in input
	if err := dec.Decode(&dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		// If the JSON contains a field which cannot be mapped to the target destination
		// then Decode() will now return an error message in the format "json: unknown
		// field "<name>"". We check for this, extract the field name from the error,
		// and interpolate it into our custom error message.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		// If the request body exceeds 1MB in size the decode will now fail with the
		// error "http: request body too large".
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err := dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func (app *Application) failedValidationResponse(w http.ResponseWriter, errors map[string]string) {
	problem := models.ErrorProblem{
		Title:  "input data not valid",
		Status: http.StatusUnprocessableEntity,
		Detail: ErrContentNotValid,
		Errors: errors,
	}
	if err := app.writeError(w, http.StatusUnprocessableEntity, problem, nil); err != nil {
		app.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Application) badRequestResponse(w http.ResponseWriter, err error) {
	problem := models.ErrorProblem{
		Title:  "input data not valid",
		Status: http.StatusBadRequest,
		Detail: err.Error(),
	}
	if err = app.writeError(w, http.StatusBadRequest, problem, nil); err != nil {
		app.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
