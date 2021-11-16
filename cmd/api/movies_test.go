package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"yamda_go/internal/config"
	"yamda_go/internal/data/provider"
	"yamda_go/internal/services"

	"github.com/julienschmidt/httprouter"
	is2 "github.com/matryer/is"
)

var app *Application = nil

func setupTestCase(_ *testing.T) func(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	cfg, _ := config.New("./../../debug.env")
	app = &Application{
		log:      logger,
		config:   cfg,
		movieSvc: services.New(provider.New(cfg)),
	}

	//TODO check if data is seeded, if not do it and clean afterwards
	return func(t *testing.T) {
		//some teardown
		app = nil
	}
}

func TestApplication_CreateMovieHandler_NoInputReceived(t *testing.T) {
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
	expectedBody := `{"title":"input data not valid","status":400,"detail":"body must not be empty"}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_InputContainsUnknownFields(t *testing.T) {
	is := is2.New(t)

	teardown := setupTestCase(t)
	defer teardown(t)

	content := `{"title": "Moana", "rating":"PG"}` //rating field is unknown to our api
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusBadRequest, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":400,"detail":"body contains unknown key \"rating\""}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_InputContainsMultipleMovies(t *testing.T) {
	is := is2.New(t)

	teardown := setupTestCase(t)
	defer teardown(t)

	content := `{"title": "Moana"}{"title": "Top Gun"}` //two movies
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusBadRequest, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":400,"detail":"body must only contain a single JSON value"}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_InputContainsGarbageContent(t *testing.T) {
	is := is2.New(t)

	teardown := setupTestCase(t)
	defer teardown(t)

	content := `{"title": "Moana", "runtime": "125 mins", "year": 2020, "genres":["drama"]} :-))`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusBadRequest, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":400,"detail":"body must only contain a single JSON value"}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_InputNotValidJSON(t *testing.T) {
	is := is2.New(t)

	teardown := setupTestCase(t)
	defer teardown(t)

	content := make([]byte, 1234)
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(string(content)))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusBadRequest, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":400,"detail":"body contains badly-formed JSON (at character 1)"}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_TitleIsEmptyOrLongerThan500Bytes(t *testing.T) {
	is := is2.New(t)

	teardown := setupTestCase(t)
	defer teardown(t)

	content := `{"title": " ", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":422,"detail":"content of movie entity is not valid","errors":{"title":"must not be empty or more than 500 bytes long"}}`
	is.Equal(expectedBody, string(body))

	content = `{"title": "tbXgdREwqSjfnDiDHUDadZPWHXPxFrzquhjpNLBjMXBnydPiwfXPxFrzquhjpNLBjMXBnydPiwfXPxFrzquhjpNLBjMXBnydPiwfXPxFrzquhjpNLBjMXBnydPiwfyuMhBTdaFRHJGDVLFkwCTvTGRcEFqNtkfhUTiAnYzQRXaRtaRrGKaJSbncPpjDAZBWtcCkWzZvJDaMgRzYBQNSpGShDhLmfUcrCRMvjpxZRSNWtqUyVHuQXVwvXKdbtzkYaWGiLgeBxNYwZgMjLtuMWbedRAvjYSWNvtzBzDAvShPWixdaFvWiMmhpVzbmZQQWEJJRaxwDBvYMDSKDWqjreFQfEBUaKrmBufecwWmEcjWmzBtKckqRddWMKacRHdNMutCBtjjTZbkbbhGvpFetxpDXZcHQBBHiWVZEHGDawmJwAntwQHtErEFvbANcrbUJhanuykDhYktjrdkuFRmQVdPFnWcRmrbKpkLtaNCcubDEuyRQYarRyjSQXWFBXbQELUPJRLBMgmNdwUdcmAXaTkwiyzdAdURrhcSCScUCZNzHGTwjWmwSpNAvRAAkuLfb", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama"]}`
	req = httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w = httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)
	resp = w.Result()
	body, _ = io.ReadAll(resp.Body)

	expectedBody = `{"title":"input data not valid","status":422,"detail":"content of movie entity is not valid","errors":{"title":"must not be empty or more than 500 bytes long"}}`
	is.Equal(expectedBody, string(body))
	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
}

func TestApplication_CreateMovieHandler_YearIsEmptyOrHasInvalidRange(t *testing.T) {
	is := is2.New(t)
	teardown := setupTestCase(t)
	defer teardown(t)

	content := `{"title": "Casablanca", "runtime": "125 mins", "year": 1800, "genres": ["historical","drama"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":422,"detail":"content of movie entity is not valid","errors":{"year":"must be greater than 1888"}}`
	is.Equal(expectedBody, string(body))

	content = `{"title": "Casablanca", "runtime": "125 mins", "genres": ["historical","drama"], "year": 2030}`
	req = httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w = httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)
	resp = w.Result()
	body, _ = io.ReadAll(resp.Body)
	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody = `{"title":"input data not valid","status":422,"detail":"content of movie entity is not valid","errors":{"year":"must not be in the future"}}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_RuntimeIsNegativeInteger(t *testing.T) {
	is := is2.New(t)
	teardown := setupTestCase(t)
	defer teardown(t)

	content := `{"title": "Casablanca", "runtime": "-1 mins", "year": 2020, "genres": ["historical","drama"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":422,"detail":"content of movie entity is not valid","errors":{"runtime":"must be a positive integer"}}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_GenresIsEmpty(t *testing.T) {
	is := is2.New(t)
	teardown := setupTestCase(t)
	defer teardown(t)

	content := `{"title": "Casablanca", "runtime": "125 mins", "year": 2020, "genres": []}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":422,"detail":"content of movie entity is not valid","errors":{"genres":"must contain at least 1 genre"}}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_GenresMustNotExceed5(t *testing.T) {
	is := is2.New(t)
	teardown := setupTestCase(t)
	defer teardown(t)

	content := `{"title": "Casablanca", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama","spy","fiction","romance","fantasy"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":422,"detail":"content of movie entity is not valid","errors":{"genres":"must not contain more than 5 genres"}}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_GenresMustBeUnique(t *testing.T) {
	is := is2.New(t)
	teardown := setupTestCase(t)
	defer teardown(t)

	content := `{"title": "Casablanca", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama","historical"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":422,"detail":"content of movie entity is not valid","errors":{"genres":"must not contain duplicate values"}}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_GetMovieHandler_Ok(t *testing.T) {
	is := is2.New(t)

	teardown := setupTestCase(t)
	defer teardown(t)

	req := httptest.NewRequest("GET", "localhost:8081/v1/movies/1", nil)
	w := httptest.NewRecorder()
	p := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: "1",
		},
	}
	app.GetMovieHandler(w, req, p)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusOK, resp.StatusCode)
	is.Equal("application/json", resp.Header.Get("Content-Type"))

	movie := `{"movie":{"id":1,"title":"The Last Samurai","runtime":"824637739352 mins","genres":["drama"," history"],"year":2015,"version":1}}`
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

func TestApplication_GetMovieHandler_MovieNotFound(t *testing.T) {
	is := is2.New(t)

	teardown := setupTestCase(t)
	defer teardown(t)

	req := httptest.NewRequest("GET", "localhost:8081/v1/movies/700", nil)
	w := httptest.NewRecorder()
	p := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: "700",
		},
	}
	app.GetMovieHandler(w, req, p)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusNotFound, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))

	movie := `{"title":"movie not found","status":404,"detail":"movie with id 700 not found"}`
	is.Equal(movie, string(body))
}
