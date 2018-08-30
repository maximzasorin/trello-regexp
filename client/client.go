package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Client allows to make API requests
type Client interface {
	GetMe() (member *Member, err error)
}

// NewClient Create new client
func NewClient(httpClient *http.Client) Client {
	return &client{httpClient}
}

type client struct {
	HttpClient *http.Client
}

// Member represents API member
type Member struct {
	ID string
}

// GetMe return current member
func (c client) GetMe() (*Member, error) {
	var m Member
	err := c.makeGetRequest("members/me", &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (c client) makeGetRequest(resource string, v interface{}) error {
	// Get member ID
	res, err := c.HttpClient.Get("https://trello.com/1/" + resource)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Parse response
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}
