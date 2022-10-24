package validator

import (
	"testing"

	is2 "github.com/matryer/is"
)

func TestIsValid_True(t *testing.T) {
	is := is2.New(t)
	validator := New()
	//validator with no errors must be valid
	is.True(validator.IsValid())
}

func TestIsValid_False(t *testing.T) {
	is := is2.New(t)
	validator := New()
	validator.AddError("key", "error message")
	is.True(validator.IsValid() == false)
}

func TestIsValid_False_UsingCheck(t *testing.T) {
	is := is2.New(t)
	validator := New()
	predicate := false
	validator.Check(predicate, "key", "error message")
	is.True(validator.IsValid() == false)
}

func TestAddError_AddsError(t *testing.T) {
	is2 := is2.New(t)
	validator := New()
	validator.AddError("key", "error message")
	is2.Equal(len(validator.Errors), 1)
}

func TestAddError_AddsError_UsingCheck(t *testing.T) {
	is2 := is2.New(t)
	validator := New()
	predicate := false
	validator.Check(predicate, "key", "error message")
	is2.Equal(len(validator.Errors), 1)
}
