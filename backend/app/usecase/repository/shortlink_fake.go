package repository

import (
	"errors"
	"time"

	"github.com/short-d/short/backend/app/entity"
)

var _ ShortLink = (*ShortLinkFake)(nil)

// ShortLinkFake accesses ShortLink information in short_link table through SQL.
type ShortLinkFake struct {
	shortLinks map[string]entity.ShortLink
}

// IsAliasExist checks whether a given alias exist in short_link table.
func (s ShortLinkFake) IsAliasExist(alias string) (bool, error) {
	_, ok := s.shortLinks[alias]
	return ok, nil
}

// CreateShortLink inserts a new ShortLink into short_link table.
func (s *ShortLinkFake) CreateShortLink(shortLink entity.ShortLink) error {
	isExist, err := s.IsAliasExist(shortLink.Alias)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("alias exists")
	}
	s.shortLinks[shortLink.Alias] = shortLink
	return nil
}

// GetShortLinkByAlias finds an ShortLink in short_link table given alias.
func (s ShortLinkFake) GetShortLinkByAlias(alias string) (entity.ShortLink, error) {
	isExist, err := s.IsAliasExist(alias)
	if err != nil {
		return entity.ShortLink{}, err
	}
	if !isExist {
		return entity.ShortLink{}, errors.New("alias not found")
	}
	shortLink := s.shortLinks[alias]
	return shortLink, nil
}

// GetShortLinksByAliases finds all ShortLink for a list of aliases
func (s ShortLinkFake) GetShortLinksByAliases(aliases []string) ([]entity.ShortLink, error) {
	if len(aliases) == 0 {
		return []entity.ShortLink{}, nil
	}

	var shortLinks []entity.ShortLink
	for _, alias := range aliases {
		shortLink, err := s.GetShortLinkByAlias(alias)

		if err != nil {
			return shortLinks, err
		}
		shortLinks = append(shortLinks, shortLink)
	}
	return shortLinks, nil
}

// UpdateShortLink updates an existing ShortLink with new properties.
func (s ShortLinkFake) UpdateShortLink(oldAlias string, newShortLink entity.ShortLink) (entity.ShortLink, error) {
	prevShortLink, ok := s.shortLinks[oldAlias]
	if !ok {
		return entity.ShortLink{}, errors.New("alias not found")
	}

	now := time.Now().UTC()
	createdBy := prevShortLink.CreatedBy
	createdAt := prevShortLink.CreatedAt
	return entity.ShortLink{
		Alias:     newShortLink.Alias,
		LongLink:  newShortLink.LongLink,
		ExpireAt:  newShortLink.ExpireAt,
		CreatedBy: createdBy,
		CreatedAt: createdAt,
		UpdatedAt: &now,
	}, nil
}

// NewShortLinkFake creates in memory ShortLink repository
func NewShortLinkFake(shortLinks map[string]entity.ShortLink) ShortLinkFake {
	return ShortLinkFake{
		shortLinks: shortLinks,
	}
}
