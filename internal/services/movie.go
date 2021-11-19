package services

import (
	"yamda_go/internal/data/provider"
	"yamda_go/internal/models"
)

type IMovieService interface {
	Insert(*models.Movie) (*models.Movie, error)
	Get(int64) (*models.Movie, error)
	Update(models.Movie) error
	Delete(int64) error
}

type MovieService struct {
	p provider.IMovieProvider
}

func New(p provider.IMovieProvider) IMovieService {
	return &MovieService{
		p,
	}
}

func (s *MovieService) Insert(m *models.Movie) (*models.Movie, error) {
	return s.p.Insert(m)
}

func (s *MovieService) Get(id int64) (*models.Movie, error) {
	m, err := s.p.Get(id)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MovieService) Update(m models.Movie) error {
	return s.p.Update(m)
}

func (s *MovieService) Delete(id int64) error {
	return s.p.Delete(id)
}
