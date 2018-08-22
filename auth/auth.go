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
	GetCallbackHandler(func(store.MemberAccessToken)) http.HandlerFunc
}

// Config represents config for auth
type Config struct {
	Name        string
	CallbackURL string
	Key         string
	Secret      string
}

// NewAuth create new auth object
func NewAuth(config *Config) Auth {
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

	return auth{config, consumer, tokens}
}

type auth struct {
	config   *Config
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

func (a auth) GetCallbackHandler(callback func(store.MemberAccessToken)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		code := values.Get("oauth_verifier")
		token := values.Get("oauth_token")

		requestToken, ok := a.tokens[token]

		if !ok {
			http.Error(w, "Cannot find request token", http.StatusInternalServerError)
		}

		accessToken, err := a.consumer.AuthorizeToken(requestToken, code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		callback(store.MemberAccessToken(*accessToken))
	}
}
