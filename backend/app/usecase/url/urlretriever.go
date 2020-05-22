package url

import (
	"fmt"
	"time"

	"github.com/short-d/short/backend/app/entity"
	"github.com/short-d/short/backend/app/usecase/repository"
)

var _ Retriever = (*RetrieverPersist)(nil)

// Retriever represents ShortLink retriever
type Retriever interface {
	GetURL(alias string, expiringAt *time.Time) (entity.ShortLink, error)
	GetURLsByUser(user entity.User) ([]entity.ShortLink, error)
}

// RetrieverPersist represents ShortLink retriever that fetches ShortLink from persistent
// storage, such as database
type RetrieverPersist struct {
	urlRepo             repository.URL
	userURLRelationRepo repository.UserURLRelation
}

// GetURL retrieves ShortLink from persistent storage given alias
func (r RetrieverPersist) GetURL(alias string, expiringAt *time.Time) (entity.ShortLink, error) {
	if expiringAt == nil {
		return r.getURL(alias)
	}
	return r.getURLExpireAfter(alias, *expiringAt)
}

func (r RetrieverPersist) getURLExpireAfter(alias string, expiringAt time.Time) (entity.ShortLink, error) {
	url, err := r.getURL(alias)
	if err != nil {
		return entity.ShortLink{}, err
	}

	if url.ExpireAt == nil {
		return url, nil
	}

	if expiringAt.After(*url.ExpireAt) {
		return entity.ShortLink{}, fmt.Errorf("url expired (alias=%s,expiringAt=%v)", alias, expiringAt)
	}

	return url, nil
}

func (r RetrieverPersist) getURL(alias string) (entity.ShortLink, error) {
	url, err := r.urlRepo.GetByAlias(alias)
	if err != nil {
		return entity.ShortLink{}, err
	}

	return url, nil
}

// GetURLsByUser retrieves URLs created by given user from persistent storage
func (r RetrieverPersist) GetURLsByUser(user entity.User) ([]entity.ShortLink, error) {
	aliases, err := r.userURLRelationRepo.FindAliasesByUser(user)
	if err != nil {
		return []entity.ShortLink{}, err
	}

	return r.urlRepo.GetByAliases(aliases)
}

// NewRetrieverPersist creates persistent ShortLink retriever
func NewRetrieverPersist(urlRepo repository.URL, userURLRelationRepo repository.UserURLRelation) RetrieverPersist {
	return RetrieverPersist{
		urlRepo:             urlRepo,
		userURLRelationRepo: userURLRelationRepo,
	}
}
