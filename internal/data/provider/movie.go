package provider

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
	"yamda_go/internal/config"
	"yamda_go/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type IMovieProvider interface {
	GetMovie(id int64) (*models.Movie, error)
	CreateMovie(*models.Movie) (*models.Movie, error)
	UpdateMovie(models.Movie) error
	DeleteMovie(id int64) error
}

type MovieProvider struct {
	db *sql.DB
}

func New(set *config.Settings) IMovieProvider {
	db, err := sql.Open(set.DriverName, set.ConnString)
	if err != nil {
		log.Fatal(err)
	}
	//validate connection to database is open correctly
	if err = db.Ping(); err != nil {
		log.Fatal(err.Error())
	}
	db.SetConnMaxLifetime(time.Minute * time.Duration(set.ConnMaxLifetime))
	db.SetMaxOpenConns(set.ConnMaxOpen)
	db.SetMaxIdleConns(set.ConnMaxIdle)
	return &MovieProvider{
		db: db,
	}
}

func (p *MovieProvider) GetMovie(id int64) (*models.Movie, error) {
	query := "SELECT * FROM Movie WHERE Id=?"
	stmt, err := p.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	//use it to scan data from row
	tmp := struct {
		ID      int64
		Title   string
		Runtime int32
		Genres  string
		Year    int32
		Version int
	}{}

	if err = stmt.QueryRow(id).Scan(&tmp.ID, &tmp.Title, &tmp.Runtime, &tmp.Genres, &tmp.Year, &tmp.Version); err != nil {
		return nil, fmt.Errorf("error scanning data from DB into internal struct: %s", err)
	}

	//build the movie model correctly
	m := models.Movie{
		ID:        tmp.ID,
		Title:     tmp.Title,
		Runtime:   models.Runtime(tmp.Runtime),
		Genres:    strings.Split(tmp.Genres, ","),
		Year:      tmp.Year,
		Version:   tmp.Version,
		CreatedAt: time.Now(), //todo change to use row field
	}

	return &m, nil
}

func (p *MovieProvider) CreateMovie(m *models.Movie) (*models.Movie, error) {
	query := `
		INSERT INTO Movie (title, runtime, genres, year, version)
		VALUES (?, ?, ?, ?, ?)
		RETURNING ID, created_at, version`
	stmtIns, err := p.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmtIns.Close()
	args := []interface{}{m.Title, m.Runtime, strings.Join(m.Genres, ", "), m.Year, m.Version}
	err = stmtIns.QueryRow(args).Scan(&m.ID, &m.CreatedAt, &m.Version)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (p *MovieProvider) UpdateMovie(m models.Movie) error {
	query := "UPDATE Movie  SET title=?, runtime=?, genres=?, year=?, version=? WHERE id=?;"
	stmtIns, err := p.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(m.Title, m.Runtime, strings.Join(m.Genres, ", "), m.Year, m.Version)
	if err != nil {
		return err
	}
	return nil
}

func (p *MovieProvider) DeleteMovie(id int64) error {
	query := "DELETE FROM Movie WHERE id=?;"
	stmtIns, err := p.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
