package store

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/coreos/bbolt"
)

// Store for everything
type Store interface {
	SaveMember(ID string, member *Member)
}

// Member represents Trello user
type Member struct {
	AccessToken MemberAccessToken
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

// SaveMember saves member to store
func (s *store) SaveMember(ID string, member *Member) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("members"))
		if err != nil {
			return err
		}

		buf, err := json.Marshal(&member)
		if err != nil {
			return err
		}

		bucket.Put([]byte(ID), buf)

		return nil
	})

	if err != nil {
		log.Fatal("Cannot save token for member #" + ID)
	}

	fmt.Println("Member #" + ID + " saved")
}
