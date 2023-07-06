//go:build !skipset1
// +build !skipset1

package provider

import (
	"log"
	"testing"
	"yamda_go/internal/models"

	is2 "github.com/matryer/is"
)

var user *models.User

func setupPreviouslyInsertedUser() func() {
	user = &models.User{
		Name:      "Slim",
		Email:     "slim@example.com",
		Password:  *models.NewPassword("", []byte("aab4162e5cb0c0da025002c99ef526db")),
		Activated: false,
	}

	prov := NewUserProvider(envConfigs, logger)
	res, err := prov.Insert(user)
	if err != nil {
		log.Fatal(err)
	}
	user = res

	return func() {
		//some teardown
		_ = prov.Delete(user.ID)
		user = nil
	}
}

func TestProvider_InsertUser_Ok(t *testing.T) {
	is := is2.New(t)

	user = &models.User{
		Name:      "Slim",
		Email:     "slim@example.com",
		Activated: false,
	}
	user.Password = *models.NewPassword("", []byte{})
	_ = user.Password.Set("aab4162e5cb0c0da025002c99ef526db")

	prov := NewUserProvider(envConfigs, logger)

	res, err := prov.Insert(user)
	is.NoErr(err)

	is.Equal(res.Version, 1)
	is.Equal(res.Name, "Slim")
	is.Equal(res.Email, "slim@example.com")
	isMatch, err := res.Password.Matches("aab4162e5cb0c0da025002c99ef526db")
	is.NoErr(err)
	is.Equal(isMatch, true)
	is.True(res.ID != 0)

	//clean it up
	err = prov.Delete(res.ID)
	is.NoErr(err)
}

func TestProvider_InsertUser_DuplicateEmail_NotOK(t *testing.T) {
	is := is2.New(t)

	teardown := setupPreviouslyInsertedUser()
	defer teardown()

	user = &models.User{
		Name:      "Slim",
		Email:     "slim@example.com",
		Activated: false,
	}
	user.Password = *models.NewPassword("", []byte{})
	_ = user.Password.Set("aab4162e5cb0c0da025002c99ef526db")

	prov := NewUserProvider(envConfigs, logger)

	//insert two times
	res, err := prov.Insert(user)

	is.Equal(res, nil)
	is.True(err != nil)
}
