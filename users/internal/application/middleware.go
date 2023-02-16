package application

import (
	"ecom-users/internal/repository"
	"ecom-users/internal/validator"
	"errors"
	"net/http"
	"strings"
)

func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		// w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, repository.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if ValidateTokenPlainText(v, token); !v.Valid() {
			app.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		// =====================
		tokenModel, err := app.Models.Tokens.Get(token, repository.ScopeAuthentication)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				app.InvalidAuthenticationTokenResponse(w, r)
				return
			}
		}

		
		user, err := app.Models.Users.GetByID(tokenModel.UserID)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				app.InvalidAuthenticationTokenResponse(w, r)
				return
			}
		}
		// =====================

		r = app.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}