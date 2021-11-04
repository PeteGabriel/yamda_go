package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"yamda_go/internal/config"
)

//TODO generate this automatically at build time
const (
	version = "0.0.1"
	ID = "id"
)

//Application type contains all dependencies for the top layer of
//the API.
type Application struct {
	log *log.Logger
	config *config.Settings
}


//ParseId parses the parameter id present in a given
//route params sent via parameters to this function.
//If the parameter id is not found an error will be returned.
//method doesn't use any dependencies from our application struct
//it could just be a regular function, rather than a method on application.
//But in general, I suggest setting up all your application-specific handlers
//and helpers so that they are methods on application.
//It helps maintain consistency in your code structure,
//and also future-proofs your code for when those handlers and helpers change later,
//and they do need access to a dependency.
func (app *Application) ParseId(p httprouter.Params) (int64, error) {
	num := p.ByName(ID)
	id, err := strconv.ParseInt(num, 10, 64)
	if err != nil || id < 1 {
		return -1, errors.New("invalid id parameter from route parameters") }
	return id, nil
}

//envelope type. Allows to insert types inside and self-document them in JSON responses.
type envelope map[string]interface{}


//Helper method for sending JSON responses. This takes the destination
// http.ResponseWriter, the HTTP status code to send, the data to encode to JSON, and a
// header map containing any additional HTTP headers we want to include in the response.
//By default, the header "Content-Type" is set to "application/json".
func (app *Application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	resp, err := json.Marshal(data)
	if err != nil {
		app.log.Println(err)
		return errors.Wrap(err, "an error happened at the server level")
	}
	//apply all values to respective header keys
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)//status must be the last write
	w.Write(resp)
	return nil
}