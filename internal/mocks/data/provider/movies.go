package provider

import (
	"yamda_go/internal/data"
	"yamda_go/internal/models"
)

type MovieProviderMock struct {
	GetMovieMock     func(int64) (*models.Movie, error)
	GetAllMoviesMock func(data.Search) ([]*models.Movie, *models.Metadata, error)
	CreateMovieMock  func(*models.Movie) (*models.Movie, error)
	UpdateMovieMock  func(models.Movie) error
	DeleteMovieMock  func(int64) error
}

func (m MovieProviderMock) Get(id int64) (*models.Movie, error) {
	return m.GetMovieMock(id)
}

func (m MovieProviderMock) GetAll(params data.Search) ([]*models.Movie, *models.Metadata, error) {
	return m.GetAllMoviesMock(params)
}

func (m MovieProviderMock) Insert(movie *models.Movie) (*models.Movie, error) {
	return m.CreateMovieMock(movie)
}

func (m MovieProviderMock) Update(movie models.Movie) error {
	return m.UpdateMovieMock(movie)
}

func (m MovieProviderMock) Delete(id int64) error {
	return m.DeleteMovieMock(id)
}
