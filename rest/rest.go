package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/maximzasorin/trello-regexp/auth"
)

// Rest represents api
type Rest interface {
	GetAuthMiddleware() mux.MiddlewareFunc
}

// NewRest creates new instance of rest
func NewRest(jwt auth.Jwt) Rest {
	return &rest{jwt}
}

type rest struct {
	jwt auth.Jwt
}

func (rst *rest) GetAuthMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rst.jwt.VerifyTokenFromHeader(r)

			next.ServeHTTP(w, r)
		})
	}
}
