package store

import (
	"fmt"

	"github.com/coreos/bbolt"
	"github.com/mrjones/oauth"
)

// Store for everything
type Store interface {
	SaveToken(*MemberAccessToken)
}

// Member represents Trello user
type Member struct {
	ID          string
	accessToken oauth.AccessToken
}

// MemberAccessToken represents user's token for access to Trello API
type MemberAccessToken struct {
	Token          string
	Secret         string
	AdditionalData map[string]string
}

// NewStore creates new store
func NewStore(db *bolt.DB) Store {
	return &store{db}
}

type store struct {
	db *bolt.DB
}

func (s *store) SaveToken(accessToken *MemberAccessToken) {
	fmt.Println("token saved " + accessToken.Token)
}
