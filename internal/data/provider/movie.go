package provider

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
	"yamda_go/internal/config"
	"yamda_go/internal/data"
	"yamda_go/internal/jsonlog"
	"yamda_go/internal/models"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type IMovieProvider interface {
	Get(int64) (*models.Movie, error)
	GetAll(data.Search) ([]*models.Movie, *models.Metadata, error)
	Insert(*models.Movie) (*models.Movie, error)
	Update(models.Movie) error
	Delete(int64) error
}

type MovieProvider struct {
	db      *sql.DB
	configs *config.Settings
}

func NewMovieProvider(set *config.Settings, log *jsonlog.Logger) IMovieProvider {
	db, err := sql.Open(set.DriverName, set.ConnString)
	if err != nil {
		log.PrintFatal(err, nil)
	}

	//validate connection to database is open correctly
	if err = db.Ping(); err != nil {
		log.PrintFatal(err, nil)
	}

	db.SetConnMaxLifetime(time.Minute * time.Duration(set.ConnMaxLifetime))
	db.SetMaxOpenConns(set.ConnMaxOpen)
	db.SetMaxIdleConns(set.ConnMaxIdle)

	return &MovieProvider{
		db:      db,
		configs: set,
	}
}

func (p *MovieProvider) Get(id int64) (*models.Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.configs.HttpReqTimeout)*time.Second)
	defer cancel()

	query := "SELECT sleep(10), Id, created_at, title, year, runtime, genres, version FROM Movie WHERE Id=?"
	stmt, err := p.db.Prepare(query)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	defer stmt.Close()

	//drivers like MariaDB have different behaviors
	row := stmt.QueryRowContext(ctx, id)
	if (*row).Err() != nil {
		return nil, ErrRecordNotFound
	}

	//use it to scan data from row
	tmp := struct {
		sleep    []byte
		ID       int64
		CreateAt []uint8
		Title    string
		Year     int32
		Runtime  int32
		Genres   string
		Version  int
	}{}

	if err = row.Scan(&tmp.sleep, &tmp.ID, &tmp.CreateAt, &tmp.Title, &tmp.Year, &tmp.Runtime, &tmp.Genres, &tmp.Version); err != nil {
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
		CreatedAt: tmp.CreateAt,
	}

	return &m, nil
}

func (p *MovieProvider) Insert(m *models.Movie) (*models.Movie, error) {
	query := `
		INSERT INTO Movie (title, runtime, genres, year, version)
		VALUES (?, ?, ?, ?, ?)
		RETURNING ID, created_at, version`
	stmtIns, err := p.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmtIns.Close()
	args := []interface{}{m.Title, 157, strings.Join(m.Genres, ", "), m.Year, m.Version}
	err = stmtIns.QueryRow(args...).Scan(&m.ID, &m.CreatedAt, &m.Version)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (p *MovieProvider) Update(m models.Movie) error {
	query := "UPDATE Movie  SET title=?, runtime=?, genres=?, year=?, version=? WHERE id=? AND version=?;"
	stmtIns, err := p.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(m.Title, m.Runtime, strings.Join(m.Genres, ", "), m.Year, m.Version+1, m.ID, m.Version)
	if err != nil {
		return err
	}
	return nil
}

func (p *MovieProvider) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := "DELETE FROM Movie WHERE id = ?"
	res, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return ErrRecordNotFound
	}
	return nil
}

func (p *MovieProvider) GetAll(params data.Search) ([]*models.Movie, *models.Metadata, error) {
	query := `
      SELECT COUNT(*) OVER(), Id, created_at, title, year, runtime, genres, version
      FROM Movie
      WHERE (LOWER(title) like LOWER(?)) AND (genres LIKE ?)
      ORDER BY %s %s, id ASC
      LIMIT %d OFFSET %d;`

	query = fmt.Sprintf(query,
		params.Filters.GetSortColumn(),
		params.Filters.GetSortDirection(),
		params.Filters.GetPageSize(),
		params.Filters.GetPageOffset())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//format params to allow a contains inside query
	title, genres := "%"+params.Title+"%", "%"+params.Genres+"%"
	rows, err := p.db.QueryContext(ctx, query, title, genres)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	//goland:noinspection GoPreferNilSlice
	movies := []*models.Movie{}

	totalRecords := 0
	for rows.Next() {

		m := struct {
			ID        int64
			CreatedAt []uint8
			Title     string
			Year      int32
			Runtime   int32
			Genres    string
			Version   int
		}{}
		if err = rows.Scan(
			&totalRecords,
			&m.ID,
			&m.CreatedAt,
			&m.Title,
			&m.Year,
			&m.Runtime,
			&m.Genres,
			&m.Version); err != nil {
			return nil, nil, fmt.Errorf("error scanning data from DB into internal struct: %s", err)
		}

		//quick solution. We'll check it at a later stage
		movie := models.Movie{
			ID:        m.ID,
			Title:     m.Title,
			Runtime:   models.Runtime(m.Runtime),
			Genres:    strings.Split(m.Genres, ","),
			Year:      m.Year,
			Version:   m.Version,
			CreatedAt: m.CreatedAt,
		}
		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	meta := models.New(totalRecords, params.Filters.Page, params.Filters.GetPageSize())
	return movies, &meta, nil
}
