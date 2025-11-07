package user

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/eugenius-watchman/ecom_go_rest_api/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()


	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUsers(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func scanRowsIntoUsers(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	var user types.User

	err := s.db.QueryRow("SELECT id, firstName, lastName, email, password, createdAt FROM users WHERE id = ?", id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName, 
        &user.Email,
        &user.Password,
        &user.CreatedAt,
	)
	if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }
    
    return &user, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password, createdAt) VALUES (?, ?, ?, ?, ?)", 
	user.FirstName, user.LastName, user.Email, user.Password, time.Now())
	if err != nil {
		return err
	}

	return nil
}