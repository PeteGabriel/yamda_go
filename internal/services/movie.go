package services

import (
	"yamda_go/internal/data/provider"
	"yamda_go/internal/models"
)

type IMovieService interface {
	CreateMovie(*models.Movie) (*models.Movie, error)
	GetMovie(int64) (*models.Movie, error)
	UpdateMovie(models.Movie) error
	DeleteMovie(int64) error
}

type MovieService struct {
	p provider.IMovieProvider
}

func New(p provider.IMovieProvider) IMovieService {
	return &MovieService{
		p,
	}
}

func (s *MovieService) CreateMovie(m *models.Movie) (*models.Movie, error) {
	return s.p.CreateMovie(m)
}

func (s *MovieService) GetMovie(id int64) (*models.Movie, error) {
	m, err := s.p.GetMovie(id)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MovieService) UpdateMovie(m models.Movie) error {
	return s.p.UpdateMovie(m)
}

func (s *MovieService) DeleteMovie(id int64) error {
	return s.p.DeleteMovie(id)
}
