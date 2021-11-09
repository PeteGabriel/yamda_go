package services

import (
	"yamda_go/internal/data/provider"
	"yamda_go/internal/models"
)

type IMovieService interface {
	CreateMovie(movie models.Movie) (bool, error)
	GetMovie(id int64) (*models.Movie, error)
}

type MovieService struct {
	p provider.IMovieProvider
}

func (s *MovieService) CreateMovie(m models.Movie) (bool, error) {
	return false, nil
}

func (s *MovieService) GetMovie(id int64) (*models.Movie, error) {
	m, err := s.p.GetMovie(id)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func New(p provider.IMovieProvider) IMovieService {
	return &MovieService{
		p,
	}
}