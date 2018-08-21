package auth

import (
	"log"
	"net/http"

	"github.com/maximzasorin/trello-regexp/store"
	"github.com/mrjones/oauth"
)

// Auth allows auth with trello
type Auth interface {
	GetRedirectHandler() http.HandlerFunc
	GetCallbackHandler() http.HandlerFunc
}

// Config represents config for auth
type Config struct {
	Name        string
	CallbackURL string
	Key         string
	Secret      string
}

// NewAuth create new auth object
func NewAuth(config *Config, store store.Store) Auth {
	consumer := oauth.NewConsumer(
		config.Key,
		config.Secret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://trello.com/1/OAuthGetRequestToken",
			AuthorizeTokenUrl: "https://trello.com/1/OAuthAuthorizeToken",
			AccessTokenUrl:    "https://trello.com/1/OAuthGetAccessToken",
		},
	)

	consumer.AdditionalAuthorizationUrlParams["name"] = config.Name
	consumer.AdditionalAuthorizationUrlParams["expiration"] = "never"
	consumer.AdditionalAuthorizationUrlParams["scope"] = "read,write"

	tokens := make(map[string]*oauth.RequestToken)

	return auth{config, store, consumer, tokens}
}

type auth struct {
	config   *Config
	store    store.Store
	consumer *oauth.Consumer
	tokens   map[string]*oauth.RequestToken
}

func (a auth) GetRedirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, requestURL, err := a.consumer.GetRequestTokenAndUrl(a.config.CallbackURL)
		if err != nil {
			log.Fatal(err)
		}
		a.tokens[token.Token] = token
		http.Redirect(w, r, requestURL, http.StatusTemporaryRedirect)
	}
}

func (a auth) GetCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		code := values.Get("oauth_verifier")
		token := values.Get("oauth_token")

		accessToken, err := a.consumer.AuthorizeToken(a.tokens[token], code)
		if err != nil {
			log.Fatal(err)
		}

		a.store.SaveToken(accessToken.Token)
	}
}
