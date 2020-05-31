package resolver

import (
	"time"

	"github.com/short-d/short/backend/app/adapter/gqlapi/scalar"
	"github.com/short-d/short/backend/app/entity"
	"github.com/short-d/short/backend/app/usecase/authenticator"
	"github.com/short-d/short/backend/app/usecase/changelog"
	"github.com/short-d/short/backend/app/usecase/shortlink"
)

// AuthQuery represents GraphQL query resolver that acts differently based
// on the identify of the user
type AuthQuery struct {
	authToken          *string
	authenticator      authenticator.Authenticator
	changeLog          changelog.ChangeLog
	shortLinkRetriever shortlink.Retriever
}

// URLArgs represents possible parameters for ShortLink endpoint
type URLArgs struct {
	Alias       string
	ExpireAfter *scalar.Time
}

// ShortLink retrieves an ShortLink persistent storage given alias and expiration time.
func (v AuthQuery) URL(args *URLArgs) (*URL, error) {
	var expireAt *time.Time
	if args.ExpireAfter != nil {
		expireAt = &args.ExpireAfter.Time
	}

	s, err := v.shortLinkRetriever.GetShortLink(args.Alias, expireAt)
	if err != nil {
		return nil, err
	}
	return &URL{url: s}, nil
}

// ChangeLog retrieves full ChangeLog from persistent storage
func (v AuthQuery) ChangeLog() (ChangeLog, error) {
	user, err := viewer(v.authToken, v.authenticator)
	if err != nil {
		return newChangeLog([]entity.Change{}, nil), ErrInvalidAuthToken{}
	}

	changeLog, err := v.changeLog.GetChangeLog()
	if err != nil {
		return ChangeLog{}, err
	}

	lastViewedAt, err := v.changeLog.GetLastViewedAt(user)
	return newChangeLog(changeLog, lastViewedAt), err
}

// URLs retrieves urls created by a given user from persistent storage
func (v AuthQuery) URLs() ([]URL, error) {
	user, err := viewer(v.authToken, v.authenticator)
	if err != nil {
		return []URL{}, ErrInvalidAuthToken{}
	}

	shortLinks, err := v.shortLinkRetriever.GetShortLinksByUser(user)
	if err != nil {
		return []URL{}, err
	}

	var gqlURLs []URL
	for _, v := range shortLinks {
		gqlURLs = append(gqlURLs, newURL(v))
	}

	return gqlURLs, nil
}

func newAuthQuery(
	authToken *string,
	authenticator authenticator.Authenticator,
	changeLog changelog.ChangeLog,
	shortLinkRetriever shortlink.Retriever,
) AuthQuery {
	return AuthQuery{
		authToken:          authToken,
		authenticator:      authenticator,
		changeLog:          changeLog,
		shortLinkRetriever: shortLinkRetriever,
	}
}
