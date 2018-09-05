package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/maximzasorin/trello-regexp/auth"
	"github.com/maximzasorin/trello-regexp/store"
)

// Rest represents api
type Rest interface {
	GetAuthMiddleware() mux.MiddlewareFunc
	Expose(router *mux.Router)
}

// NewRest creates new instance of rest
func NewRest(jwt auth.Jwt, st store.Store) Rest {
	return &rest{jwt, st}
}

type rest struct {
	jwt   auth.Jwt
	store store.Store
}

func (rst *rest) Expose(router *mux.Router) {
	router.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{}")
	})
}

func (rst *rest) GetAuthMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			memberID, err := rst.jwt.VerifyTokenFromHeader(r)
			if err != nil {
				http.Error(w, "Invalid token.", http.StatusUnauthorized)
				return
			}

			member, err := rst.store.GetMember(memberID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if member == nil {
				http.Error(w, "Unknown member.", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
