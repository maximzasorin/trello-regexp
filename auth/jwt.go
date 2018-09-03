package auth

import (
	"net/http"
	"time"

	webTokens "github.com/dgrijalva/jwt-go"
	"github.com/maximzasorin/trello-regexp/store"
)

// Jwt allow works with JWT
type Jwt interface {
	AuthMember(w http.ResponseWriter, member *store.Member) error
	GetSecret() string
}

// NewJwt creates new jwt instance
func NewJwt(secret string, cookieName string) Jwt {
	return &jwt{secret, cookieName}
}

type jwt struct {
	secret     string
	cookieName string
}

// AuthMember pass jwt token to cookies
func (j *jwt) AuthMember(w http.ResponseWriter, member *store.Member) error {
	token := webTokens.NewWithClaims(webTokens.SigningMethodHS256, webTokens.MapClaims{
		"member_id": member.ID,
	})

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return err
	}

	cookieVar := http.Cookie{
		Name:    j.cookieName,
		Value:   tokenString,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 24),
	}

	http.SetCookie(w, &cookieVar)

	return nil
}

// GetSecret Returns secret
func (j *jwt) GetSecret() string {
	return j.secret
}
