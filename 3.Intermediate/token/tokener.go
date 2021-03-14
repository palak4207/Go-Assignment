package token

import (
	"time"
)

// Tokener ...
type Tokener interface {
	CreateToken(username string, duration time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}
