package data

import (
	"strings"
	"yamda_go/internal/validator"
)

type Search struct {
	Title   string
	Genres  string
	Filters Filter
}

type Filter struct {
	Page         int
	PageSize     int
	Sort         string   //field to be used while sorting results
	SortSafelist []string //specified by each endpoint to customize their search
}

//Validate performs a few checks of the fields of a given instance of Filter.
func (f Filter) Validate(v *validator.Validator) {
	v.Check(f.Page > 0, "page", "must be bigger than 0")
	v.Check(f.Page < 10_000_000, "page", "must be smaller than 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be bigger than 0")
	v.Check(f.PageSize <= 100, "page", "must be smaller than 100")

	// Check that the sort parameter matches a value in the safelist.
	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

//GetSortColumn returns which column should be used to sort results.
func (f Filter) GetSortColumn() string {
	for _, s := range f.SortSafelist {
		if f.Sort == s {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	//should never happen since we have a validator routine for these cases.
	panic("unsafe sort parameter: " + f.Sort)
}

//GetSortDirection returns which direction (ASC, DESC) should be used while sorting results.
func (f Filter) GetSortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filter) GetPageSize() int {
	return f.PageSize
}

func (f Filter) GetPageOffset() int {
	return (f.Page - 1) * f.PageSize
}
