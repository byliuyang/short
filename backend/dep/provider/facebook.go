package provider

import (
	"github.com/short-d/app/fw"
	"github.com/short-d/short/app/adapter/facebook"
)

// FacebookClientID represents client ID used for Facebook OAuth.
type FacebookClientID string

// FacebookClientSecret represents client secret used for Facebook OAuth.
type FacebookClientSecret string

// FacebookRedirectURI represents redirect URL for facebook single sign on.
type FacebookRedirectURI string

// NewFacebookIdentityProvider creates a new Facebook OAuth client with
// FacebookClientID and FacebookClientSecret to uniquely identify clientID and
// clientSecret during dependency injection.
func NewFacebookIdentityProvider(
	req fw.HTTPRequest,
	clientID FacebookClientID,
	clientSecret FacebookClientSecret,
	redirectURI FacebookRedirectURI,
) facebook.IdentityProvider {
	return facebook.NewIdentityProvider(req, string(clientID), string(clientSecret), string(redirectURI))
}
