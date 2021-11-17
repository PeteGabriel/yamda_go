package provider

import "yamda_go/internal/models"

type MovieProviderMock struct {
	GetMovieMock    func(int64) (*models.Movie, error)
	CreateMovieMock func(models.Movie) (bool, error)
}

func (m MovieProviderMock) GetMovie(id int64) (*models.Movie, error) {
	return m.GetMovieMock(id)
}

func (m MovieProviderMock) CreateMovie(movie models.Movie) (bool, error) {
	return m.CreateMovieMock(movie)
}
