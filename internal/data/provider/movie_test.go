package provider

import (
	"testing"
	"yamda_go/internal/config"
	"yamda_go/internal/models"

	is2 "github.com/matryer/is"
)

func setupTestCase() func() {

	return func() {
		//some teardown

	}
}

func TestProvider_InsertMovie_Ok(t *testing.T) {
	is := is2.New(t)
	mov := &models.Movie{
		Title:   "Once upon a time",
		Runtime: 157,
		Year:    1998,
		Genres:  []string{"fantasy", "drama"},
		Version: 1,
	}
	cfg, _ := config.New("./../../../debug.env")
	prov := New(cfg)

	res, err := prov.Insert(mov)
	is.NoErr(err)

	is.Equal("Once upon a time", res.Title)
	is.True(157 == res.Runtime)
	is.True(1998 == res.Year)
	is.Equal([]string{"fantasy", "drama"}, res.Genres)
	is.True(res.ID > 0)
	is.Equal(1, res.Version)
	//TODO Fix is.True(res.CreatedAt < time.Now())
}
