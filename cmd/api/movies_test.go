package main

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	is2 "github.com/matryer/is"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"yamda_go/internal/config"
	"yamda_go/internal/data/provider"
	"yamda_go/internal/jsonlog"
	provmock "yamda_go/internal/mocks/data/provider"
	"yamda_go/internal/models"
)

var appMoviesTest *Application = nil

func setupTestCase(p provmock.MovieProviderMock) func() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	cfg, _ := config.New("./../../debug.env")
	appMoviesTest = &Application{
		logger:        logger,
		config:        cfg,
		movieProvider: p,
	}
	return func() {
		//some teardown
		appMoviesTest = nil
	}
}

/*********************************************************
** CREATE
*********************************************************/

func TestApplication_CreateMovieHandler_NoInputReceived(t *testing.T) {
	is := is2.New(t)

	//setup mock for movieProvider
	mock := provmock.MovieProviderMock{}
	mock.CreateMovieMock = func(movie *models.Movie) (*models.Movie, error) {
		return nil, nil //dummy return in this case
	}
	teardown := setupTestCase(mock)
	defer teardown()

	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", nil)
	w := httptest.NewRecorder()
	appMoviesTest.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusBadRequest, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"bad request","status":400,"detail":"body must not be empty"}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_InputContainsUnknownFields(t *testing.T) {
	is := is2.New(t)

	mock := provmock.MovieProviderMock{}
	mock.CreateMovieMock = func(movie *models.Movie) (*models.Movie, error) {
		return nil, nil //dummy return in this case
	}
	teardown := setupTestCase(mock)
	defer teardown()

	content := `{"title": "Moana", "runtime": "125 mins", "year": 2020, "genres":["drama"], "rating":"PG"}` //rating field is unknown to our api
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	appMoviesTest.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusBadRequest, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"bad request","status":400,"detail":"body contains unknown key \"rating\""}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_InputContainsMultipleMovies(t *testing.T) {
	is := is2.New(t)

	mock := provmock.MovieProviderMock{}
	mock.CreateMovieMock = func(movie *models.Movie) (*models.Movie, error) {
		return nil, nil //dummy return in this case
	}
	teardown := setupTestCase(mock)
	defer teardown()

	content := `{"title": "Moana"}{"title": "Top Gun"}` //two movies
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	appMoviesTest.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusBadRequest, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"bad request","status":400,"detail":"body must only contain a single JSON value"}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_InputContainsGarbageContent(t *testing.T) {
	is := is2.New(t)

	mock := provmock.MovieProviderMock{}
	mock.CreateMovieMock = func(movie *models.Movie) (*models.Movie, error) {
		return nil, nil //dummy return in this case
	}
	teardown := setupTestCase(mock)
	defer teardown()

	content := `{"title": "Moana", "runtime": "125 mins", "year": 2020, "genres":["drama"]} :-))`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	appMoviesTest.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusBadRequest, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"bad request","status":400,"detail":"body must only contain a single JSON value"}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_InputNotValidJSON(t *testing.T) {
	is := is2.New(t)

	mock := provmock.MovieProviderMock{}
	mock.CreateMovieMock = func(movie *models.Movie) (*models.Movie, error) {
		return nil, nil //dummy return in this case
	}
	teardown := setupTestCase(mock)
	defer teardown()

	content := make([]byte, 1234)
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(string(content)))
	w := httptest.NewRecorder()
	appMoviesTest.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusBadRequest, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"bad request","status":400,"detail":"body contains badly-formed JSON (at character 1)"}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_TitleIsEmptyOrLongerThan500Bytes(t *testing.T) {
	is := is2.New(t)

	mock := provmock.MovieProviderMock{}
	mock.CreateMovieMock = func(movie *models.Movie) (*models.Movie, error) {
		return nil, nil //dummy return in this case
	}
	teardown := setupTestCase(mock)
	defer teardown()

	content := `{"title": " ", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	appMoviesTest.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input validations failed","status":422,"detail":"","errors":{"title":"must not be empty or more than 500 bytes long"}}`
	is.Equal(expectedBody, string(body))

	content = `{"title": "tbXgdREwqSjfnDiDHUDadZPWHXPxFrzquhjpNLBjMXBnydPiwfXPxFrzquhjpNLBjMXBnydPiwfXPxFrzquhjpNLBjMXBnydPiwfXPxFrzquhjpNLBjMXBnydPiwfyuMhBTdaFRHJGDVLFkwCTvTGRcEFqNtkfhUTiAnYzQRXaRtaRrGKaJSbncPpjDAZBWtcCkWzZvJDaMgRzYBQNSpGShDhLmfUcrCRMvjpxZRSNWtqUyVHuQXVwvXKdbtzkYaWGiLgeBxNYwZgMjLtuMWbedRAvjYSWNvtzBzDAvShPWixdaFvWiMmhpVzbmZQQWEJJRaxwDBvYMDSKDWqjreFQfEBUaKrmBufecwWmEcjWmzBtKckqRddWMKacRHdNMutCBtjjTZbkbbhGvpFetxpDXZcHQBBHiWVZEHGDawmJwAntwQHtErEFvbANcrbUJhanuykDhYktjrdkuFRmQVdPFnWcRmrbKpkLtaNCcubDEuyRQYarRyjSQXWFBXbQELUPJRLBMgmNdwUdcmAXaTkwiyzdAdURrhcSCScUCZNzHGTwjWmwSpNAvRAAkuLfb", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama"]}`
	req = httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w = httptest.NewRecorder()
	appMoviesTest.CreateMovieHandler(w, req, nil)
	resp = w.Result()
	body, _ = io.ReadAll(resp.Body)

	expectedBody = `{"title":"input validations failed","status":422,"detail":"","errors":{"title":"must not be empty or more than 500 bytes long"}}`
	is.Equal(expectedBody, string(body))
	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
}

func TestApplication_CreateMovieHandler_YearIsEmptyOrHasInvalidRange(t *testing.T) {
	is := is2.New(t)
	mock := provmock.MovieProviderMock{}
	mock.CreateMovieMock = func(movie *models.Movie) (*models.Movie, error) {
		return nil, nil //dummy return in this case
	}
	teardown := setupTestCase(mock)
	defer teardown()

	content := `{"title": "Casablanca", "runtime": "125 mins", "year": 1800, "genres": ["historical","drama"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	appMoviesTest.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input validations failed","status":422,"detail":"","errors":{"year":"must be greater than 1888"}}`
	is.Equal(expectedBody, string(body))

	content = `{"title": "Casablanca", "runtime": "125 mins", "genres": ["historical","drama"], "year": 2030}`
	req = httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w = httptest.NewRecorder()
	appMoviesTest.CreateMovieHandler(w, req, nil)
	resp = w.Result()
	body, _ = io.ReadAll(resp.Body)
	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody = `{"title":"input validations failed","status":422,"detail":"","errors":{"year":"must not be in the future"}}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_RuntimeIsNegativeInteger(t *testing.T) {
	is := is2.New(t)
	mock := provmock.MovieProviderMock{}
	mock.CreateMovieMock = func(movie *models.Movie) (*models.Movie, error) {
		return nil, nil //dummy return in this case
	}
	teardown := setupTestCase(mock)
	defer teardown()

	content := `{"title": "Casablanca", "runtime": "-1 mins", "year": 2020, "genres": ["historical","drama"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	appMoviesTest.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input validations failed","status":422,"detail":"","errors":{"runtime":"must be a positive integer"}}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_GenresIsEmpty(t *testing.T) {
	is := is2.New(t)
	mock := provmock.MovieProviderMock{}
	mock.CreateMovieMock = func(movie *models.Movie) (*models.Movie, error) {
		return nil, nil //dummy return in this case
	}
	teardown := setupTestCase(mock)
	defer teardown()

	content := `{"title": "Casablanca", "runtime": "125 mins", "year": 2020, "genres": []}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	appMoviesTest.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input validations failed","status":422,"detail":"","errors":{"genres":"must contain at least 1 genre"}}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_GenresMustNotExceed5(t *testing.T) {
	is := is2.New(t)
	mock := provmock.MovieProviderMock{}
	mock.CreateMovieMock = func(movie *models.Movie) (*models.Movie, error) {
		return nil, nil //dummy return in this case
	}
	teardown := setupTestCase(mock)
	defer teardown()

	content := `{"title": "Casablanca", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama","spy","fiction","romance","fantasy"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	appMoviesTest.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input validations failed","status":422,"detail":"","errors":{"genres":"must not contain more than 5 genres"}}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_GenresMustBeUnique(t *testing.T) {
	is := is2.New(t)
	mock := provmock.MovieProviderMock{}
	mock.CreateMovieMock = func(movie *models.Movie) (*models.Movie, error) {
		return nil, nil //dummy return in this case
	}
	teardown := setupTestCase(mock)
	defer teardown()

	content := `{"title": "Casablanca", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama","historical"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	appMoviesTest.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"input validations failed","status":422,"detail":"","errors":{"genres":"must not contain duplicate values"}}`
	is.Equal(expectedBody, string(body))
}

func TestApplication_CreateMovieHandler_Ok(t *testing.T) {
	is := is2.New(t)
	mock := provmock.MovieProviderMock{}
	mock.CreateMovieMock = func(movie *models.Movie) (*models.Movie, error) {
		movie.ID = 12
		movie.Version = 1
		return movie, nil
	}
	teardown := setupTestCase(mock)
	defer teardown()

	content := `{"title": "Casablanca", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama"]}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	appMoviesTest.CreateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusCreated, resp.StatusCode)
	is.Equal("application/json", resp.Header.Get("Content-Type"))
	is.Equal("/v1/movies/12", resp.Header.Get("Location"))
	expectedBody := `{"movie":{"id":12,"title":"Casablanca","runtime":"125 mins","genres":["historical","drama"],"year":2020,"version":1}}`
	is.Equal(expectedBody, string(body))
}

/*********************************************************
** GET MOVIE
*********************************************************/

func TestApplication_GetMovieHandler_Ok(t *testing.T) {
	is := is2.New(t)

	mock := provmock.MovieProviderMock{}
	mock.GetMovieMock = func(id int64) (*models.Movie, error) {
		return &models.Movie{
			ID:      1,
			Title:   "The Last Samurai",
			Runtime: 127,
			Genres:  []string{"drama", " history"},
			Year:    2015,
			Version: 1,
		}, nil
	}
	teardown := setupTestCase(mock)
	defer teardown()

	req := httptest.NewRequest("GET", "localhost:8081/v1/movies/1", nil)
	w := httptest.NewRecorder()
	p := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: "1",
		},
	}
	appMoviesTest.GetMovieHandler(w, req, p)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusOK, resp.StatusCode)
	is.Equal("application/json", resp.Header.Get("Content-Type"))

	movie := `{"movie":{"id":1,"title":"The Last Samurai","runtime":"127 mins","genres":["drama"," history"],"year":2015,"version":1}}`
	is.Equal(movie, string(body))
}

func TestApplication_GetMovieHandler_BadMovieId(t *testing.T) {
	is := is2.New(t)

	mock := provmock.MovieProviderMock{}
	mock.GetMovieMock = func(id int64) (*models.Movie, error) {
		return nil, errors.New("")
	}
	teardown := setupTestCase(mock)
	defer teardown()

	req := httptest.NewRequest("GET", "localhost:8081/v1/movies/7p", nil)
	w := httptest.NewRecorder()
	p := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: "7p",
		},
	}
	appMoviesTest.GetMovieHandler(w, req, p)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusBadRequest, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))

	movie := `{"title":"bad request","status":400,"detail":"invalid id parameter from route parameters"}`
	is.Equal(movie, string(body))
}

func TestApplication_GetMovieHandler_MovieNotFound(t *testing.T) {
	is := is2.New(t)

	mock := provmock.MovieProviderMock{}
	mock.GetMovieMock = func(id int64) (*models.Movie, error) {
		return nil, provider.ErrRecordNotFound
	}
	teardown := setupTestCase(mock)
	defer teardown()

	req := httptest.NewRequest("GET", "localhost:8081/v1/movies/700", nil)
	w := httptest.NewRecorder()
	p := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: "700",
		},
	}
	appMoviesTest.GetMovieHandler(w, req, p)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusNotFound, resp.StatusCode)
	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))

	expectedBody := `{"title":"resource not found","status":404,"detail":"movie with id 700 not found"}`
	is.Equal(expectedBody, string(body))
}

/*********************************************************
** UPDATE
*********************************************************/

func TestApplication_UpdateMovieHandler_WithoutSpecifyingID(t *testing.T) {
	is := is2.New(t)
	mock := provmock.MovieProviderMock{}
	mock.UpdateMovieMock = func(m models.Movie) error {
		return nil
	}
	teardown := setupTestCase(mock)
	defer teardown()

	content := `{"title": "Casablanca", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama"]}`
	req := httptest.NewRequest("PATCH", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	appMoviesTest.UpdateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	expectedBody := `{"title":"input validations failed","status":422,"detail":"","errors":{"ID":"ID must be provided and bigger than 0"}}`
	is.Equal(string(body), expectedBody)
}

func TestApplication_UpdateMovieHandler_Ok(t *testing.T) {
	is := is2.New(t)
	mock := provmock.MovieProviderMock{}
	mock.UpdateMovieMock = func(m models.Movie) error {
		return nil
	}
	teardown := setupTestCase(mock)
	defer teardown()

	content := `{"id": 1, "title": "Casablanca", "runtime": "125 mins", "year": 2020, "genres": ["historical","drama"]}`
	req := httptest.NewRequest("PATCH", "localhost:8081/v1/movies", strings.NewReader(content))
	w := httptest.NewRecorder()
	appMoviesTest.UpdateMovieHandler(w, req, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	is.Equal(http.StatusNoContent, resp.StatusCode)
	is.Equal(string(body), "")
}

func TestApplication_UpdateMovieHandler_UnknownField(t *testing.T) {
	_ = is2.New(t)
}

/*********************************************************
** Partial Update
*********************************************************/

func TestApplication_UpdateMovieHandler_Partially(t *testing.T) {
	is := is2.New(t)
	mock := provmock.MovieProviderMock{}
	mock.UpdateMovieMock = func(m models.Movie) error {
		return nil
	}
	//this simulates what existed previously on the database
	mock.GetMovieMock = func(id int64) (*models.Movie, error) {
		return &models.Movie{
			ID:      id,
			Title:   "Let's grow old together",
			Runtime: 100,
			Year:    2000,
			Genres:  []string{"romance"},
		}, nil
	}
	teardown := setupTestCase(mock)
	defer teardown()

	content := `{"runtime": "125 mins", "genres": ["drama"]}`
	req := httptest.NewRequest("PATCH", "localhost:8081/v1/movies/1", strings.NewReader(content))
	w := httptest.NewRecorder()
	appMoviesTest.PartialUpdateMovieHandler(w, req, httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: "1",
		}})

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	is.Equal(http.StatusOK, resp.StatusCode)
	expectedBody := `{"movie":{"id":1,"title":"Let's grow old together","runtime":"125 mins","genres":["drama"],"year":2000,"version":0}}`
	is.Equal(string(body), expectedBody)
}

/*********************************************************
** DELETE
*********************************************************/
func TestApplication_DeleteMovieHandler_Ok(t *testing.T) {
	is := is2.New(t)
	mock := provmock.MovieProviderMock{}
	mock.DeleteMovieMock = func(i int64) error {
		return nil
	}
	teardown := setupTestCase(mock)
	defer teardown()

	req := httptest.NewRequest("DELETE", "localhost:8081/v1/movies/1", nil)
	w := httptest.NewRecorder()
	p := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: "1",
		},
	}
	appMoviesTest.DeleteMovieHandler(w, req, p)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	is.Equal(http.StatusOK, resp.StatusCode)
	is.Equal(string(body), "")
}

func TestApplication_DeleteMovieHandler_MovieNotFound(t *testing.T) {
	is := is2.New(t)
	mock := provmock.MovieProviderMock{}
	mock.DeleteMovieMock = func(i int64) error {
		return provider.ErrRecordNotFound
	}
	teardown := setupTestCase(mock)
	defer teardown()

	req := httptest.NewRequest("DELETE", "localhost:8081/v1/movies/1", nil)
	w := httptest.NewRecorder()
	p := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: "1",
		},
	}
	appMoviesTest.DeleteMovieHandler(w, req, p)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	is.Equal(http.StatusNotFound, resp.StatusCode)

	is.Equal("application/problem+json", resp.Header.Get("Content-Type"))
	expectedBody := `{"title":"resource not found","status":404,"detail":"movie with id 1 not found"}`
	is.Equal(expectedBody, string(body))
}
