package provider

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
	"yamda_go/internal/config"
	"yamda_go/internal/models"
)

type IMovieProvider interface {
	GetMovie(id int64) (*models.Movie, error)
	CreateMovie(*models.Movie) (bool, error)
}

type MovieProvider struct {
	db *sql.DB
}

func New(set *config.Settings) IMovieProvider{
	db, err := sql.Open(set.DriverName, set.ConnString)
	if err != nil {
		//todo handle this
		log.Fatal(err)
	}
	err = db.Ping() //validate connection to database is open correctly
	if err != nil {
		log.Fatal(err.Error()) // proper error handling instead of panic in your app
	}

	db.SetConnMaxLifetime(time.Minute * time.Duration(set.ConnMaxLifetime))
	db.SetMaxOpenConns(set.ConnMaxOpen)
	db.SetMaxIdleConns(set.ConnMaxIdle)
	var provider = &MovieProvider{
		db: db,
	}

	return provider
}

func (p *MovieProvider) GetMovie(id int64) (*models.Movie, error) {
	query := "SELECT * FROM Movie WHERE Id=?"
	stmt, err := p.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	prod := models.Movie{
		ID:        0,
		Title:     "",
		Runtime:   0,
		Genres:    nil,
		Year:      0,
		Version:   0,
		CreatedAt: time.Time{},
	}
	err = stmt.QueryRow(id).Scan(
		&prod.ID,
		&prod.Title,
		&prod.Runtime,
		&prod.Genres,
		&prod.Year,
		&prod.Version,
		&prod.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &prod, nil
}

func (p *MovieProvider) CreateMovie(*models.Movie) (bool, error) {

	return false, nil
}