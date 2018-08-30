package auth

import (
	"net/http"

	"github.com/maximzasorin/trello-regexp/client"
	"github.com/maximzasorin/trello-regexp/store"
	"github.com/mrjones/oauth"
)

// Auth allows auth with trello
type Auth interface {
	GetRedirectHandler() http.HandlerFunc
	GetCallbackHandler() http.HandlerFunc
	GetHttpClient(*oauth.AccessToken) (*http.Client, error)
}

// Config represents config for auth
type Config struct {
	Name         string
	CallbackURL  string
	TrelloKey    string
	TrelloSecret string
}

// NewAuth create new auth object
func NewAuth(config *Config, store store.Store) Auth {
	consumer := oauth.NewConsumer(
		config.TrelloKey,
		config.TrelloSecret,
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

	return &auth{config, store, consumer, tokens}
}

type auth struct {
	config   *Config
	store    store.Store
	consumer *oauth.Consumer
	tokens   map[string]*oauth.RequestToken
}

// GetRedirectHandler returns http handler for redirect to Oauth provider
func (a *auth) GetRedirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, requestURL, err := a.consumer.GetRequestTokenAndUrl(a.config.CallbackURL)
		if err != nil {
			a.triggerServerError(w, err.Error())
		}
		a.tokens[token.Token] = token
		http.Redirect(w, r, requestURL, http.StatusTemporaryRedirect)
	}
}

// GetRedirectHandler returns http handler for process callback from Oauth provider
func (a *auth) GetCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		code := values.Get("oauth_verifier")
		token := values.Get("oauth_token")

		requestToken, ok := a.tokens[token]
		if !ok {
			a.triggerServerError(w, "Cannot find request token")
			return
		}

		// Get AccessToken
		accessToken, err := a.consumer.AuthorizeToken(requestToken, code)
		if err != nil {
			a.triggerServerError(w, err.Error())
			return
		}

		// Create Client for API requests
		httpClient, err := a.GetHttpClient(accessToken)
		if err != nil {
			a.triggerServerError(w, err.Error())
			return
		}

		client := client.NewClient(httpClient)
		clientMember, err := client.GetMe()
		if err != nil {
			a.triggerServerError(w, err.Error())
			return
		}

		// Update member data
		member, err := a.store.GetMember(clientMember.ID)
		if err != nil {
			a.triggerServerError(w, err.Error())
			return
		}

		// Update access token
		if member == nil {
			member = &store.Member{
				ID: clientMember.ID,
			}
		}
		member.AccessToken = *accessToken

		err = a.store.SaveMember(member)
		if err != nil {
			a.triggerServerError(w, err.Error())
		}
	}
}

// GetHttpClient return HTTP client for make API responses
func (a *auth) GetHttpClient(accessToken *oauth.AccessToken) (*http.Client, error) {
	return a.consumer.MakeHttpClient(accessToken)
}

func (a *auth) triggerServerError(w http.ResponseWriter, err string) {
	http.Error(w, err, http.StatusInternalServerError)
}
