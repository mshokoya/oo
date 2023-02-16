package application

import (
	"ecom-users/internal/repository"
	"ecom-users/internal/validator"
)

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be valid email address")
} 

func ValidatePasswordPlainText(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 characters long")
	v.Check(len(password) <= 30, "password", "must be at most 30 characters long")
}

func ValidateUsername(v *validator.Validator, username string) {
	v.Check(username != "", "username", "must be provided")
	v.Check(len(username) >= 5, "username", "must be at least 5 characters long")
	v.Check(len(username) <= 12, "username", "must be at most 12 characters long")
}

func ValidateUser(v *validator.Validator, user *repository.User) {
	v.Check(user.Firstname != "", "firstname", "must be provided")
	v.Check(len(user.Firstname) <= 30, "firstname", "must be less than 30 characters long")
	ValidateUsername(v, user.Username)
	ValidateEmail(v, user.Email)
	ValidatePasswordPlainText(v, *user.Password.Plaintext)

	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}