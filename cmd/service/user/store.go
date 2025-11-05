package user

import (
	"database/sql"

	"github.com/eugenius-watchman/ecom_go_rest_api/types"
)


type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error)
