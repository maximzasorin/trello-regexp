package store

import (
	"encoding/json"

	"github.com/coreos/bbolt"
	"github.com/mrjones/oauth"
)

// Store for everything
type Store interface {
	SaveMember(member *Member) error
}

// Member represents Trello user
type Member struct {
	ID          string
	AccessToken oauth.AccessToken
}

// NewStore creates new store
func NewStore(db *bolt.DB) Store {
	return &store{db}
}

type store struct {
	db *bolt.DB
}

// SaveMember saves member to store
func (s *store) SaveMember(member *Member) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("members"))
		if err != nil {
			return err
		}

		buf, err := json.Marshal(&member)
		if err != nil {
			return err
		}

		bucket.Put([]byte(member.ID), buf)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
