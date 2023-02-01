package provider

import (
	is2 "github.com/matryer/is"
	"os"
	"testing"
	"yamda_go/internal/config"
	"yamda_go/internal/jsonlog"
	"yamda_go/internal/models"
)

var user *models.User

func TestProvider_InsertUser_Ok(t *testing.T) {
	is := is2.New(t)

	user = &models.User{
		Name:  "Slim",
		Email: "slim@example.com",
		Password: models.Password{
			Hash: []byte("aab4162e5cb0c0da025002c99ef526db"),
		},
		Activated: false,
	}
	cfg, _ := config.New("./../../../debug.env")
	prov := NewUserProvider(cfg, jsonlog.New(os.Stdout, jsonlog.LevelInfo))

	res, err := prov.Insert(user)
	is.NoErr(err)

	is.Equal(res.Version, 1)
	is.Equal(res.Name, "Slim")
	is.Equal(res.Email, "slim@example.com")
	is.Equal(res.Password.Hash, []byte("aab4162e5cb0c0da025002c99ef526db"))
	is.True(res.ID != 0)

	//clean it up
	err = prov.Delete(res.ID)
	is.NoErr(err)
}
