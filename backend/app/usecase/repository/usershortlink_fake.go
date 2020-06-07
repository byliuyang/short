package repository

import (
	"errors"

	"github.com/short-d/short/backend/app/entity"
)

var _ UserShortLink = (*UserShortLinkFake)(nil)

// UserShortLinkFake represents in memory implementation of User-ShortLink relationship accessor.
type UserShortLinkFake struct {
	users      []entity.User
	shortLinks []entity.ShortLink
}

// CreateRelation creates many to many relationship between User and ShortLink.
func (u *UserShortLinkFake) CreateRelation(user entity.User, shortLink entity.ShortLink) error {
	if u.IsRelationExist(user, shortLink) {
		return errors.New("relationship exists")
	}
	u.users = append(u.users, user)
	u.shortLinks = append(u.shortLinks, shortLink)
	return nil
}

// FindAliasesByUser fetches the aliases of all the ShortLinks created by the given user.
func (u UserShortLinkFake) FindAliasesByUser(user entity.User) ([]string, error) {
	var aliases []string
	for idx, currUser := range u.users {
		if currUser.ID != user.ID {
			continue
		}
		aliases = append(aliases, u.shortLinks[idx].Alias)
	}
	return aliases, nil
}

// IsRelationExist checks whether the an ShortLink is own by a given user.
func (u UserShortLinkFake) IsRelationExist(user entity.User, shortLink entity.ShortLink) bool {
	for idx, currUser := range u.users {
		if currUser.ID != user.ID {
			continue
		}

		if u.shortLinks[idx].Alias == shortLink.Alias {
			return true
		}
	}
	return false
}

// NewUserShortLinkRepoFake creates UserShortLinkFake
func NewUserShortLinkRepoFake(users []entity.User, shortLinks []entity.ShortLink) UserShortLinkFake {
	return UserShortLinkFake{
		users:      users,
		shortLinks: shortLinks,
	}
}
