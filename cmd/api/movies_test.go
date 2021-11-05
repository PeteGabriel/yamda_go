package main

import (
	"github.com/julienschmidt/httprouter"
	is2 "github.com/matryer/is"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"yamda_go/internal/config"
)

var app *Application = nil

func setupTestCase(t *testing.T) func(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)
	cfg, _ := config.New("./debug.env")
	app = &Application{
		log:    logger,
		config: cfg,
	}
	return func(t *testing.T) {
		//some teardown
		app = nil
	}
}

func TestApplication_CreateMovieHandler_BadInput(t *testing.T) {
	is := is2.New(t)

	teardown := setupTestCase(t)
	defer teardown(t)

	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", nil)
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusBadRequest, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":400,"detail":"input data could not be decoded into expected structure"}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_GetMovieHandler_Ok(t *testing.T) {
	is := is2.New(t)

	teardown := setupTestCase(t)
	defer teardown(t)

	req := httptest.NewRequest("GET", "localhost:8081/v1/movies/7", nil)
	w := httptest.NewRecorder()
	p := httprouter.Params{
		httprouter.Param{
		Key:   "id",
		Value: "7",
		},
	}
	app.GetMovieHandler(w, req, p)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusOK, resp.StatusCode)
	is.Equal("application/json", resp.Header.Get("Content-Type"))

	movie := `{"movie":{"id":1,"title":"Casablanca","genres":["drama","war","romance"],"version":1}}`
	is.Equal(movie, string(body))
}

func TestApplication_GetMovieHandler_BadMovieId(t *testing.T) {
	is := is2.New(t)

	teardown := setupTestCase(t)
	defer teardown(t)

	req := httptest.NewRequest("GET", "localhost:8081/v1/movies/7p", nil)
	w := httptest.NewRecorder()
	p := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: "7p",
		},
	}
	app.GetMovieHandler(w, req, p)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusBadRequest, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))

	movie := `{"title":"error trying to parse id from route:","status":400,"detail":"invalid id parameter from route parameters"}`
	is.Equal(movie, string(body))
}