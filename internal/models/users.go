package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
	"yamda_go/internal/validator"
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type Password struct {
	//pointer to help distinguish between password not being present in versus a password which is the empty string "".
	plaintext *string
	hash      []byte
}

func NewPassword(p string, h []byte) *Password {
	return &Password{
		plaintext: &p,
		hash:      h,
	}
}

// Matches compare the given password (hash) against the saved hash
// to verify if passwords match.
func (p *Password) Matches(plaintext string) (bool, error) {
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

func (p *Password) GetHash() []byte {
	return p.hash
}

func (p *Password) SetHash(h []byte) error {
	if p.hash != nil {
		return errors.New("password already has an hash value")
	}
	p.hash = h
	return nil
}

// Set password. Generates a hash of the password and set's it.
func (p *Password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintext
	p.hash = hash
	return nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, pw string) {
	v.Check(pw != "", "password", "must be provided")
	v.Check(len(pw) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(pw) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, usr *User) {
	v.Check(usr.Name != "", "name", "must be provided")
	v.Check(len(usr.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, usr.Email)

	if usr.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *usr.Password.plaintext)
	}

	//panic due to this being a bug of our own code.
	if usr.Password.hash == nil {
		panic("missing password hash for user")
	}
}
