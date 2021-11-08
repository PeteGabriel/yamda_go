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

	"github.com/julienschmidt/httprouter"
	is2 "github.com/matryer/is"
)

var app *Application = nil

func setupTestCase(t *testing.T) func(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
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

func TestApplication_CreateMovieHandler_Created(t *testing.T) {
	is := is2.New(t)

	teardown := setupTestCase(t)
	defer teardown(t)


	content := `{"title": "Moana", "runtime": "107 mins"}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(string(content)))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	//body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusCreated, resp.StatusCode)
	is.Equal("application/json", resp.Header.Get("Content-Type"))
	//TODO is.Equal("", string(body))
}

func TestApplication_CreateMovieHandler_TitleIsEmptyOrLongerThan500Bytes(t *testing.T) {
	is := is2.New(t)

	teardown := setupTestCase(t)
	defer teardown(t)


	content := `{"title": " ", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama","historical"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(string(content)))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":400,"detail":"movie title must be non-empty"}`
	is.Equal(expectedBody, string(body))

	content = `{"title": "tbXgdREwqSjfnDiDHUDadZPWHXPxFrzquhjpNLBjMXBnydPiwfyuMhBTdaFRHJGDVLFkwCTvTGRcEFqNtkfhUTiAnYzQRXaRtaRrGKaJSbncPpjDAZBWtcCkWzZvJDaMgRzYBQNSpGShDhLmfUcrCRMvjpxZRSNWtqUyVHuQXVwvXKdbtzkYaWGiLgeBxNYwZgMjLtuMWbedRAvjYSWNvtzBzDAvShPWixdaFvWiMmhpVzbmZQQWEJJRaxwDBvYMDSKDWqjreFQfEBUaKrmBufecwWmEcjWmzBtKckqRddWMKacRHdNMutCBtjjTZbkbbhGvpFetxpDXZcHQBBHiWVZEHGDawmJwAntwQHtErEFvbANcrbUJhanuykDhYktjrdkuFRmQVdPFnWcRmrbKpkLtaNCcubDEuyRQYarRyjSQXWFBXbQELUPJRLBMgmNdwUdcmAXaTkwiyzdAdURrhcSCScUCZNzHGTwjWmwSpNAvRAAkuLfb", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama","historical"]}`
	req = httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(string(content)))
	w = httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)
	resp = w.Result()
	body, _ = io.ReadAll(resp.Body)

	expectedBody = `{"title":"input data not valid","status":422,"detail":"movie title exceeds lenght."}`
	is.Equal(expectedBody, string(body))
	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
}

func TestApplication_CreateMovieHandler_YearIsEmptyOrHasInvalidRange(t *testing.T) {
	is := is2.New(t)
	teardown := setupTestCase(t)
	defer teardown(t)

	content := `{"title": "Casablanca", "runtime": "125 mins", "year": -1, "genres": ["historical","drama","historical"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(string(content)))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":422,"detail":"movie year must be non-empty"}`
	is.Equal(expectedBody, string(body))


	content = `{"title": "Casablanca", "runtime": "125 mins", "genres": ["historical","drama","historical"], "year": 1886}`
	req = httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(string(content)))
	w = httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)
	resp = w.Result()
	body, _ = io.ReadAll(resp.Body)
	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody = `{"title":"input data not valid","status":422,"detail":"movie year must be inside range 1888 and current year"}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_RuntimeIsEmptyOrANegativeInteger(t *testing.T) {
	is := is2.New(t)
	teardown := setupTestCase(t)
	defer teardown(t)

	content := `{"title": "Casablanca", "runtime": "", "year": 2020, "genres": ["historical","drama","historical"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(string(content)))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":422,"detail":"movie runtime must be non-empty"}`
	is.Equal(expectedBody, string(body))

	content = `{"runtime": "-26 mins", "title": "Casablanca", "year": 2020, "genres": ["historical","drama","historical"]}`
	req = httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(string(content)))
	w = httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp = w.Result()
	body, _ = io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody = `{"title":"input data not valid","status":422,"detail":"movie runtime must be non-negative"}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_GenresIsEmpty(t *testing.T) {
	is := is2.New(t)
	teardown := setupTestCase(t)
	defer teardown(t)

	content := `{"title": "Casablanca", "runtime": "125 mins", "year": 2020, "genres": []}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(string(content)))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":422,"detail":"movie genres must be non-empty"}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_GenresMustNotExceed5(t *testing.T) {
	is := is2.New(t)
	teardown := setupTestCase(t)
	defer teardown(t)

	content := `{"title": "Casablanca", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama","spy","fiction","romance","fantasy"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(string(content)))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":422,"detail":"a movie must have between one and five (unique) genres"}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_GenresMustBeUnique(t *testing.T) {
	is := is2.New(t)
	teardown := setupTestCase(t)
	defer teardown(t)

	content := `{"title": "Casablanca", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama","historical"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(string(content)))
	w := httptest.NewRecorder()
	app.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input data not valid","status":422,"detail":"movie genres must be unique"}`
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
