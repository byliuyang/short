package payload

import (
	"errors"

	"github.com/short-d/app/fw"
	"github.com/short-d/short/app/entity"
)

const emailKey = "email"

var _ Payload = (*Email)(nil)

// Email represents a payload that contains user's email
type Email struct {
	user entity.User
}

// GetTokenPayload retrieves the token payload representation of the email
// payload.
func (e Email) GetTokenPayload() fw.TokenPayload {
	return map[string]interface{}{
		emailKey: e.user.Email,
	}
}

// GetUser retrieves user info represented by the email payload.
func (e Email) GetUser() entity.User {
	return e.user
}

var _ Factory = (*EmailFactory)(nil)

// EmailFactory produces email payload.
type EmailFactory struct {
}

// FromTokenPayload parses token payload into email payload.
func (e EmailFactory) FromTokenPayload(tokenPayload fw.TokenPayload) (Payload, error) {
	JSONEmail := tokenPayload[emailKey]
	email, ok := JSONEmail.(string)
	if !ok {
		return nil, errors.New("expect payload to contain email")
	}

	user := entity.User{
		Email: email,
	}
	return Email{user: user}, nil
}

// FromUser converts user info into email payload.
func (e EmailFactory) FromUser(user entity.User) (Payload, error) {
	if user.Email == "" {
		return nil, errors.New("user email cannot be empty")
	}
	return Email{user: user}, nil
}

// NewEmailFactory creates email payload factory.
func NewEmailFactory() EmailFactory {
	return EmailFactory{}
}
