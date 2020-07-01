package resolver

import (
	"errors"
	"fmt"
	"time"

	"github.com/short-d/short/backend/app/adapter/gqlapi/scalar"
	"github.com/short-d/short/backend/app/entity"
	"github.com/short-d/short/backend/app/usecase/authenticator"
	"github.com/short-d/short/backend/app/usecase/changelog"
	"github.com/short-d/short/backend/app/usecase/shortlink"
)

// AuthMutation represents GraphQL mutation resolver that acts differently based
// on the identify of the user
type AuthMutation struct {
	authToken        *string
	authenticator    authenticator.Authenticator
	changeLog        changelog.ChangeLog
	shortLinkCreator shortlink.Creator
	shortLinkUpdater shortlink.Updater
}

// ShortLinkInput represents possible ShortLink attributes
type ShortLinkInput struct {
	LongLink    *string
	CustomAlias *string
	ExpireAt    *time.Time
}

// TODO(#840): remove this business logic and move it to use cases
func (s *ShortLinkInput) isEmpty() bool {
	return *s == ShortLinkInput{}
}

// TODO(#840): remove this business logic and move it to use cases
func (s *ShortLinkInput) longLink() string {
	if s.LongLink == nil {
		return ""
	}
	return *s.LongLink
}

// TODO(#840): remove this business logic and move it to use cases
func (s *ShortLinkInput) customAlias() string {
	if s.CustomAlias == nil {
		return ""
	}
	return *s.CustomAlias
}

// TODO(#840): remove this business logic and move it to use cases
func (s *ShortLinkInput) createUpdate() *entity.ShortLink {
	if s.isEmpty() {
		return nil
	}

	return &entity.ShortLink{
		Alias:    s.customAlias(),
		LongLink: s.longLink(),
		ExpireAt: s.ExpireAt,
	}

}

// CreateShortLinkArgs represents the possible parameters for CreateShortLink endpoint
type CreateShortLinkArgs struct {
	ShortLink ShortLinkInput
	IsPublic  bool
}

// CreateShortLink creates mapping between an alias and a long link for a given user
func (a AuthMutation) CreateShortLink(args *CreateShortLinkArgs) (*ShortLink, error) {
	user, err := viewer(a.authToken, a.authenticator)
	if err != nil {
		return nil, ErrInvalidAuthToken{}
	}

	longLink := args.ShortLink.longLink()
	customAlias := args.ShortLink.CustomAlias
	u := entity.ShortLink{
		LongLink: longLink,
		ExpireAt: args.ShortLink.ExpireAt,
	}

	isPublic := args.IsPublic

	newShortLink, err := a.shortLinkCreator.CreateShortLink(u, customAlias, user, isPublic)
	if err == nil {
		return &ShortLink{shortLink: newShortLink}, nil
	}

	var (
		ae shortlink.ErrAliasExist
		l  shortlink.ErrInvalidLongLink
		c  shortlink.ErrInvalidCustomAlias
		m  shortlink.ErrMaliciousLongLink
	)
	if errors.As(err, &ae) {
		return nil, ErrAliasExist(*customAlias)
	}
	if errors.As(err, &l) {
		return nil, ErrInvalidLongLink{u.LongLink, string(l.Violation)}
	}
	if errors.As(err, &c) {
		return nil, ErrInvalidCustomAlias{*customAlias, string(c.Violation)}
	}
	if errors.As(err, &m) {
		return nil, ErrMaliciousContent(u.LongLink)
	}
	return nil, ErrUnknown{}
}

// UpdateShortLinkArgs represents the possible parameters for updateShortLink endpoint
type UpdateShortLinkArgs struct {
	OldAlias  string
	ShortLink ShortLinkInput
}

// UpdateShortLink updates the relationship between the short link and the user
func (a AuthMutation) UpdateShortLink(args *UpdateShortLinkArgs) (*ShortLink, error) {
	user, err := viewer(a.authToken, a.authenticator)
	if err != nil {
		return nil, ErrInvalidAuthToken{}
	}

	update := args.ShortLink.createUpdate()
	if update == nil {
		return nil, nil
	}

	newShortLink, err := a.shortLinkUpdater.UpdateShortLink(args.OldAlias, *update, user)
	if err != nil {
		return nil, err
	}

	return &ShortLink{shortLink: newShortLink}, nil
}

// ChangeInput represents possible properties for Change
type ChangeInput struct {
	Title           string
	SummaryMarkdown *string
}

// CreateChangeArgs represents the possible parameters for CreateChange endpoint
type CreateChangeArgs struct {
	Change ChangeInput
}

// CreateChange creates a Change in the change log
func (a AuthMutation) CreateChange(args *CreateChangeArgs) (*Change, error) {
	user, err := viewer(a.authToken, a.authenticator)
	if err != nil {
		return nil, ErrInvalidAuthToken{}
	}

	change, err := a.changeLog.CreateChange(args.Change.Title, args.Change.SummaryMarkdown, user)
	if err == nil {
		change := newChange(change)
		return &change, nil
	}

	var (
		u changelog.ErrUnauthorizedAction
	)
	if errors.As(err, &u) {
		return nil, ErrUnauthorizedAction(fmt.Sprintf("user %s is not allowed to create a change", user.ID))
	}
	return nil, ErrUnknown{}
}

// DeleteChangeArgs represents the possible parameters for DeleteChange endpoint
type DeleteChangeArgs struct {
	ID string
}

// DeleteChange removes a Change with given ID from change log
func (a AuthMutation) DeleteChange(args *DeleteChangeArgs) (*string, error) {
	user, err := viewer(a.authToken, a.authenticator)
	if err != nil {
		return nil, ErrInvalidAuthToken{}
	}

	err = a.changeLog.DeleteChange(args.ID, user)
	if err == nil {
		return &args.ID, nil
	}

	var (
		u changelog.ErrUnauthorizedAction
	)
	if errors.As(err, &u) {
		return nil, ErrUnauthorizedAction(fmt.Sprintf("user %s is not allowed to delete the change %s", user.ID, args.ID))
	}
	return nil, ErrUnknown{}
}

// UpdateChangeArgs represents the possible parameters for UpdateChange endpoint.
type UpdateChangeArgs struct {
	ID     string
	Change ChangeInput
}

// UpdateChange updates a Change with given ID in change log.
func (a AuthMutation) UpdateChange(args *UpdateChangeArgs) (*Change, error) {
	user, err := viewer(a.authToken, a.authenticator)
	if err != nil {
		return nil, ErrInvalidAuthToken{}
	}

	change, err := a.changeLog.UpdateChange(
		args.ID,
		args.Change.Title,
		args.Change.SummaryMarkdown,
		user,
	)
	if err == nil {
		change := newChange(change)
		return &change, nil
	}

	var (
		u changelog.ErrUnauthorizedAction
	)
	if errors.As(err, &u) {
		return nil, ErrUnauthorizedAction(fmt.Sprintf("user %s is not allowed to update the change %s", user.ID, args.ID))
	}
	return nil, ErrUnknown{}
}

// ViewChangeLog records the time when the user viewed the change log
func (a AuthMutation) ViewChangeLog() (scalar.Time, error) {
	user, err := viewer(a.authToken, a.authenticator)
	if err != nil {
		return scalar.Time{}, ErrInvalidAuthToken{}
	}

	lastViewedAt, err := a.changeLog.ViewChangeLog(user)
	return scalar.Time{Time: lastViewedAt}, err
}

func newAuthMutation(
	authToken *string,
	authenticator authenticator.Authenticator,
	changeLog changelog.ChangeLog,
	shortLinkCreator shortlink.Creator,
	shortLinkUpdater shortlink.Updater,
) AuthMutation {
	return AuthMutation{
		authToken:        authToken,
		authenticator:    authenticator,
		changeLog:        changeLog,
		shortLinkCreator: shortLinkCreator,
		shortLinkUpdater: shortLinkUpdater,
	}
}
