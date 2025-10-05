// Package user
package user

import "time"

type User struct {
	ID    int    `json:"id`
	Name  string `json:"name"`
	Email string `json:"email"`
	// I'm not including password as we shouldn't be returning that at all.
	CreatedAt time.Time `json:"created_at"`
}

type CreateRequest struct {
	Name     string `json:"name`
	Email    string `json:"email`
	Password string `json:"password"`
}

var (
	ErrNotFound    = Err("not found")
	ErrEmailExists = Err("email exists")
)

type Err string

func (e Err) Error() string { return string(e) }
