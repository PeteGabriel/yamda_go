package provider

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"yamda_go/internal/config"
	"yamda_go/internal/jsonlog"
	"yamda_go/internal/models"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrEmailNotFound  = errors.New("email not found")
)

type IUserProvider interface {
	Insert(user *models.User) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int64) error
}

type UserProvider struct {
	db      *sql.DB
	configs *config.Settings
}

func NewUserProvider(set *config.Settings, log *jsonlog.Logger) IUserProvider {
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

	return &UserProvider{
		db:      db,
		configs: set,
	}
}

func (u *UserProvider) Insert(user *models.User) (*models.User, error) {
	query := `INSERT INTO users (name, email, password_hash, activated)
              VALUES (?, ?, ?, ?)
              RETURNING id, created_at, version`

	insert, err := u.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer insert.Close()

	args := []interface{}{user.Name, user.Email, user.Password.GetHash(), user.Activated}
	row := insert.QueryRow(args...)
	if (*row).Err() != nil {
		return nil, ErrRecordNotFound //TODO handle with different error
	}

	if err = row.Scan(&user.ID, &user.CreatedAt, &user.Version); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserProvider) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, created_at, name, email, password_hash, activated, version
              FROM users WHERE email = ?`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.configs.HttpReqTimeout)*time.Second)
	defer cancel()

	getStatement, err := u.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer getStatement.Close()

	row := getStatement.QueryRowContext(ctx, email)
	if (*row).Err() != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrEmailNotFound
		default:
			return nil, err
		}
	}

	tmp := struct {
		ID        int64
		CreatedAt []uint8
		Name      string
		Email     string
		Password  models.Password
		Activated bool
		Version   int
	}{}

	var hash []byte
	err = row.Scan(&tmp.ID, &tmp.CreatedAt, &tmp.Name, &tmp.Email, &hash, &tmp.Activated, &tmp.Version)
	if err != nil {
		return nil, fmt.Errorf("error scanning data from DB into internal struct: %s", err)
	}

	if err := tmp.Password.SetHash(hash); err != nil {
		return nil, fmt.Errorf("error setting password hash: %s", err)
	}

	return &models.User{
		ID:        tmp.ID,
		CreatedAt: tmp.CreatedAt,
		Name:      tmp.Name,
		Email:     tmp.Email,
		Password:  tmp.Password,
		Activated: tmp.Activated,
		Version:   tmp.Version,
	}, nil
}

func (u *UserProvider) Update(user *models.User) error {
	query := `UPDATE users SET name = ?, email = ?, password_hash = ?, activated = ?, version = version + 1
             WHERE id = ? AND version = ?`
	stmtIns, err := u.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(user.Name, user.Email, user.Password.GetHash(), user.Activated, user.Version+1, user.ID, user.Version)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserProvider) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := "DELETE FROM users WHERE id = ?"
	res, err := u.db.Exec(query, id)
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
