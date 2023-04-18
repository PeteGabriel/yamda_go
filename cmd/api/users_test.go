package main

import (
	is2 "github.com/matryer/is"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"yamda_go/internal/config"
	"yamda_go/internal/jsonlog"
	provmock "yamda_go/internal/mocks/data/provider"
	"yamda_go/internal/models"
)

var appUsersTest *Application = nil

func setupUsersTestCase(p provmock.UserProviderMock) func() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	cfg, _ := config.New("./../../debug.env")
	appUsersTest = &Application{
		logger:       logger,
		config:       cfg,
		userProvider: p,
	}
	return func() {
		//some teardown
		appUsersTest = nil
	}
}

func TestApplication_RegisterNewUser_Ok(t *testing.T) {
	is := is2.New(t)

	mock := provmock.UserProviderMock{}
	mock.InsertMock = func(user *models.User) (*models.User, error) {
		return &models.User{
			ID:        1234,
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			Activated: false,
			Version:   1,
		}, nil
	}

	teardown := setupUsersTestCase(mock)
	defer teardown()

	body := `{"name": "Jason Bourne", "email": "jason@bourne.com", "password": "mysuperpwhehe"}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/users", strings.NewReader(body))
	w := httptest.NewRecorder()
	appUsersTest.RegisterUserHandler(w, req, nil)

	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusCreated, resp.StatusCode)
	expectedBody := `{"user":{"id":1234,"created_at":null,"name":"Jason Bourne","email":"jason@bourne.com","activated":false}}`
	is.Equal(expectedBody, string(respBody))
}

func TestApplication_RegisterNewUser_MissingEmail(t *testing.T) {
	is := is2.New(t)

	mock := provmock.UserProviderMock{}
	mock.InsertMock = func(user *models.User) (*models.User, error) {
		return &models.User{
			ID:        1234,
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			Activated: false,
			Version:   1,
		}, nil
	}

	teardown := setupUsersTestCase(mock)
	defer teardown()

	body := `{"name": "Jason Bourne", "password": "mysuperpwhehe"}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/users", strings.NewReader(body))
	w := httptest.NewRecorder()
	appUsersTest.RegisterUserHandler(w, req, nil)

	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	expectedBody := `{"title":"input validations failed","status":422,"detail":"","errors":{"email":"must be provided"}}`
	is.Equal(expectedBody, string(respBody))
}

func TestApplication_RegisterNewUser_MissingName(t *testing.T) {
	is := is2.New(t)

	mock := provmock.UserProviderMock{}
	mock.InsertMock = func(user *models.User) (*models.User, error) {
		return &models.User{
			ID:        1234,
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			Activated: false,
			Version:   1,
		}, nil
	}

	teardown := setupUsersTestCase(mock)
	defer teardown()

	body := `{"email": "Jason@Bourne.com", "password": "mysuperpwhehe"}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/users", strings.NewReader(body))
	w := httptest.NewRecorder()
	appUsersTest.RegisterUserHandler(w, req, nil)

	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	expectedBody := `{"title":"input validations failed","status":422,"detail":"","errors":{"name":"must be provided"}}`
	is.Equal(expectedBody, string(respBody))
}

func TestApplication_RegisterNewUser_MissingPassword(t *testing.T) {
	is := is2.New(t)

	mock := provmock.UserProviderMock{}
	mock.InsertMock = func(user *models.User) (*models.User, error) {
		return &models.User{
			ID:        1234,
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			Activated: false,
			Version:   1,
		}, nil
	}

	teardown := setupUsersTestCase(mock)
	defer teardown()

	body := `{"name": "Jason Bourne", "email": "jason@bourne.com"}`
	req := httptest.NewRequest("POST", "localhost:8081/v1/users", strings.NewReader(body))
	w := httptest.NewRecorder()
	appUsersTest.RegisterUserHandler(w, req, nil)

	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)

	is.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	expectedBody := `{"title":"input validations failed","status":422,"detail":"","errors":{"password":"must be provided"}}`
	is.Equal(expectedBody, string(respBody))
}
