package store

import (
	"fmt"

	"github.com/coreos/bbolt"
)

// Store for everything
type Store interface {
	SaveToken(token string)
}

// NewStore creates new store
func NewStore(db *bolt.DB) Store {
	return &store{db}
}

type store struct {
	db *bolt.DB
}

func (s *store) SaveToken(token string) {
	fmt.Println("token saved " + token)
}
