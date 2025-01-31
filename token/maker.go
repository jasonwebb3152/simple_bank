package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// Creates new token for specified user and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// Checks if token is valid or not.
	VerifyToken(token string) (*Payload, error)
}
