package provider

import "yamda_go/internal/models"

type MovieProviderMock struct {
	GetMovieMock    func(int64) (*models.Movie, error)
	CreateMovieMock func(*models.Movie) (*models.Movie, error)
	UpdateMovieMock func(models.Movie) error
	DeleteMovieMock func(int64) error
}

func (m MovieProviderMock) GetMovie(id int64) (*models.Movie, error) {
	return m.GetMovieMock(id)
}

func (m MovieProviderMock) CreateMovie(movie *models.Movie) (*models.Movie, error) {
	return m.CreateMovieMock(movie)
}

func (m MovieProviderMock) UpdateMovie(movie models.Movie) error {
	return m.UpdateMovieMock(movie)
}

func (m MovieProviderMock) DeleteMovie(id int64) error {
	return m.DeleteMovieMock(id)
}
