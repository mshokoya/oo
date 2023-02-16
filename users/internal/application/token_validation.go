package application

import "ecom-users/internal/validator"

func ValidateTokenPlainText(v *validator.Validator, tokenPlainText string) {
	v.Check(tokenPlainText != "", "token", "must be provided")
	v.Check(len(tokenPlainText) == 26, "token", "must be 26 bytes long")
}