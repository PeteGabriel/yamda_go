package provider

import (
	is2 "github.com/matryer/is"
	"log"
	"testing"
	"yamda_go/internal/models"
)

var mov *models.Movie

func setupTestCase() func() {
	mov = &models.Movie{
		Title:   "Once upon a time",
		Runtime: 157,
		Year:    1998,
		Genres:  []string{"fantasy", "drama"},
		Version: 1,
	}

	prov := NewMovieProvider(envConfigs, logger)
	res, err := prov.Insert(mov)
	if err != nil {
		log.Fatal(err)
	}
	mov = res

	return func() {
		//some teardown
		_ = prov.Delete(mov.ID)
		mov = nil
	}
}

func TestProvider_InsertMovie_Ok(t *testing.T) {
	is := is2.New(t)

	movie := &models.Movie{
		Title:   "Once upon a time",
		Runtime: 157,
		Year:    1998,
		Genres:  []string{"fantasy", "drama"},
		Version: 1,
	}

	prov := NewMovieProvider(envConfigs, logger)

	res, err := prov.Insert(movie)
	is.NoErr(err)

	is.Equal("Once upon a time", res.Title)
	is.True(157 == res.Runtime)
	is.True(1998 == res.Year)
	is.Equal([]string{"fantasy", "drama"}, res.Genres)
	is.True(res.ID > 0)
	is.Equal(1, res.Version)
	is.True(res.CreatedAt != nil)

	//clean it up
	err = prov.Delete(res.ID)
	is.NoErr(err)
}

func TestMovieProvider_Update_Ok(t *testing.T) {
	is := is2.New(t)

	teardown := setupTestCase()
	defer teardown()

	prov := NewMovieProvider(envConfigs, logger)

	//update year
	mov.Year = 2004
	err := prov.Update(*mov)
	is.NoErr(err)

	//assert
	tmp, err := prov.Get(mov.ID)
	is.NoErr(err)
	is.Equal(tmp.ID, mov.ID)
	is.Equal(tmp.Title, mov.Title)
	is.Equal(tmp.Runtime, mov.Runtime)
	is.True(tmp.Year == 2004)

	//clean up
	err = prov.Delete(tmp.ID)
	is.NoErr(err)
}

func TestMovieProvider_Get_Ok(t *testing.T) {
	is := is2.New(t)

	teardown := setupTestCase()
	defer teardown()

	prov := NewMovieProvider(envConfigs, logger)

	tmp, err := prov.Get(mov.ID)
	is.NoErr(err)
	is.Equal(tmp.ID, mov.ID)
	is.Equal(tmp.Title, mov.Title)
	is.Equal(tmp.Runtime, mov.Runtime)
	is.Equal(tmp.Year, mov.Year)
	is.True(tmp.CreatedAt != nil)
}
