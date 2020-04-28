package repository

import (
	"github.com/short-d/short/app/entity"
)

// URL accesses urls from storage, such as database.
type URL interface {
	IsAliasExist(alias string) (bool, error)
	GetByAlias(alias string) (entity.URL, error)
	// TODO(issue#698): change to CreateURL
	Create(url entity.URL) error
	UpdateURL(oldAlias string, newURL entity.URL) (entity.URL, error)
	GetByAliases(aliases []string) ([]entity.URL, error)
}
