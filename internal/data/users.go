package data

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type password struct {
	//pointer to help distinguish between password not being present in versus a password which is the empty string "".
	plaintext *string
	hash      []byte
}

// Matches compare the given password (hash) against the saved hash
// to verify if passwords match.
func (p *password) Matches(plaintext string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

//Set password. Generates a hash of the password and set's it.
func (p *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintext
	p.hash = hash
	return nil
}
