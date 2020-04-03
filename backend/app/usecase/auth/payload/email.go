package payload

import (
	"errors"

	"github.com/short-d/app/fw"
	"github.com/short-d/short/app/entity"
)

var _ Payload = (*Email)(nil)

type Email struct {
	User entity.User
}

func (e Email) GetTokenPayload() fw.TokenPayload {
	return map[string]interface{}{
		"email": e.User.Email,
	}
}

func (e Email) GetUser() entity.User {
	return e.User
}

var _ Factory = (*EmailFactory)(nil)

type EmailFactory struct {
}

func (e EmailFactory) FromTokenPayload(tokenPayload fw.TokenPayload) (Payload, error) {
	JSONEmail := tokenPayload["email"]
	email, ok := JSONEmail.(string)
	if !ok {
		return nil, errors.New("expect payload to contain email")
	}

	user := entity.User{
		Email: email,
	}
	return Email{User: user}, nil
}

func (e EmailFactory) FromUser(user entity.User) (Payload, error) {
	if user.Email == "" {
		return nil, errors.New("user email cannot be empty")
	}
	return Email{User: user}, nil
}

func NewEmailFactory() EmailFactory {
	return EmailFactory{}
}
