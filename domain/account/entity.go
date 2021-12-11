package account

import (
	"time"
)

const AccountSessionKeyFormat = "account:session:%s"

type AccountContextKey struct{}

// Account is a collection of proprty of account.
type Account struct {
	ID             int64      `json:"id"`
	Email          string     `json:"email"`
	Password       *string    `json:"password,omitempty"`
	FirstName      string     `json:"firstName"`
	LastName       string     `json:"lastName"`
	CreatedAt      time.Time  `json:"createdAt"`
	LastModifiedAt *time.Time `json:"lastModifiedAt"`
}