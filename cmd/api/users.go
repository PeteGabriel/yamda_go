package main

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"yamda_go/internal/data/provider"
	"yamda_go/internal/models"
	"yamda_go/internal/validator"
)

/**
* POST /v1/users -> 201 CREATED with JSON content
*
* {
*	"name": "Alice Smith",
*	"email": "alice@example.com",
*	"password": "pa55word"
* }
* Should create a new User struct containing
* these details, validate it with the ValidateUser() helper, and then pass it to our
* UserModel.Insert() method to create a new database record.
**/
func (app *Application) RegisterUserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//deserialize request. If error, return bad request response
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, err)
		return
	}

	user := &models.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  models.Password{},
		Activated: false,
	}

	if err := user.Password.Set(input.Password); err != nil {
		//our fault. Password type us ok, something wrong happened on our side
		app.serverErrorResponse(w, err)
		return
	}

	//validate data
	v := validator.New()

	if models.ValidateUser(v, user); !v.IsValid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	//try to insert data into storage
	insertedUsr, err := app.userProvider.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, provider.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, v.Errors)
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	//send response
	if err := app.writeJSON(w, http.StatusCreated, envelope{"user": insertedUsr}, nil); err != nil {
		app.serverErrorResponse(w, err)
	}
}
