package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	webTokens "github.com/dgrijalva/jwt-go"
)

// Jwt allow works with JWT
type Jwt interface {
	AuthMember(w http.ResponseWriter, memberID string) error
	VerifyTokenFromHeader(r *http.Request) (string, error)
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

// MemberClaims stores MemberID
type MemberClaims struct {
	webTokens.StandardClaims
	MemberID string
}

// AuthMember pass jwt token to cookies
func (j *jwt) AuthMember(w http.ResponseWriter, memberID string) error {
	token := webTokens.NewWithClaims(webTokens.SigningMethodHS256, &MemberClaims{
		MemberID: memberID,
	})

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    j.cookieName,
		Value:   tokenString,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 24),
	})

	return nil
}

func (j *jwt) VerifyTokenFromHeader(r *http.Request) (string, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return "", nil
	}

	token, err := webTokens.ParseWithClaims(authorization[7:], &MemberClaims{}, func(token *webTokens.Token) (interface{}, error) {
		if _, ok := token.Method.(*webTokens.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.secret), nil
	})

	if err != nil {
		return "", errors.New("Can not parse token")
	}

	claims, ok := token.Claims.(*MemberClaims)
	if !token.Valid || !ok {
		return "", errors.New("Invalid token")
	}

	return claims.MemberID, nil
}

// GetSecret Returns secret
func (j *jwt) GetSecret() string {
	return j.secret
}
