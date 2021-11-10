package models

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Runtime int32

// Define an error that our UnmarshalJSON() method can return if we're unable to parse
// or convert the JSON string successfully.
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

//UnmarshalJSON method on the Runtime type so that it satisfies the
// json.Unmarshaler interface.
func (r *Runtime) UnmarshalJSON(val []byte) error {
  // We expect that the incoming JSON value will be a string in the format 
  // "<runtime> mins", and the first thing we need to do is remove the surrounding 
	// double-quotes.
	unquotedValue, err := strconv.Unquote(string(val))
	if err != nil {
			return ErrInvalidRuntimeFormat
	}
	parts := strings.Split(unquotedValue, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}
	runtime, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	*r = Runtime(runtime)
	return nil
}

func (r *Runtime) MarshalJSON() ([]byte, error) {
	log.Println("MarshalJSON call")
	cnt := fmt.Sprintf("\"%d mins\"", r)//must be a valid JSON string
	return []byte(cnt), nil
}