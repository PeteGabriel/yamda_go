package data

import "yamda_go/internal/validator"

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
