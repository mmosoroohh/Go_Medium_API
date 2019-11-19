package middlewares

import (
	"errors"
	"github.com/mmosoroohh/Go_Medium_API/api/auth"
	"github.com/mmosoroohh/Go_Medium_API/api/responses"
	"net/http"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.ValidToken(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unathorized"))
			return
		}
		next(w, r)
	}
}