package types

import "time"

type User struct {
	// Go field name ... JSON field nam
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type RegisterUserPayload struct {
	// Go field name ... JSON field nam
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
