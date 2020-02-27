package usermodels

import "time"

// Record is a user record in the database
type Record struct {
	ID            int       `json:"id"`
	UID           string    `json:"uid"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	PasswordHash  string    `json:"-"`
	CreatedAt     time.Time `json:"created"`
	UpdatedAt     time.Time `json:"updated"`
}
