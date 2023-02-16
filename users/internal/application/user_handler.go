package application

import (
	"ecom-users/internal/repository"
	"ecom-users/internal/validator"
	"errors"
	"net/http"
)

func (app *Application) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email string `json:"email"`
		Firstname string `json:"firstname"`
	}

	err := app.ReadJSON(w, r, &input)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	user := &repository.User{
		Firstname: input.Firstname,
		Username: input.Username,
		Email: input.Email,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	
	if ValidateUser(v, user); !v.Valid() {
		app.FailedValidationResponse(w,r, v.Errors)
		return
	}

	err = app.Models.Users.Insert(*user)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateEmail):
			app.FailedValidationResponse(w,r,v.Errors)
		default:
			app.ServerErrorResponse(w,r,err)
		}
		return
	}
	
	err = app.WriteJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.ServerErrorResponse(w,r, err)
	}
}

func (app *Application) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname string `json:"lastname"`
	}

	err := app.ReadJSON(w, r, &input)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	user := repository.User{}

	if input.Username != "" {
		ValidateUsername(v, input.Username)
		user.Username = input.Username
	}

	if input.Firstname != "" {
		v.Check(len(input.Firstname) >= 2, "firstname", "Invalid firstname provided")
		user.Firstname = input.Firstname
	}

	if input.Lastname != "" {
		v.Check(len(input.Lastname) >= 2, "lastname", "Invalid lastname provided")
		user.Lastname = input.Lastname
	}

	if !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.Models.Users.Update(user)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	env := envelope{"message": "your details was successfully updated"}

	err = app.WriteJSON(w,http.StatusAccepted, env, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}