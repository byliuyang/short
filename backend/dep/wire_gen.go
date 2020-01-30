// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package dep

import (
	"database/sql"
	"time"

	"github.com/google/wire"
	"github.com/short-d/app/fw"
	"github.com/short-d/app/modern/mdcli"
	"github.com/short-d/app/modern/mddb"
	"github.com/short-d/app/modern/mdenv"
	"github.com/short-d/app/modern/mdhttp"
	"github.com/short-d/app/modern/mdio"
	"github.com/short-d/app/modern/mdlogger"
	"github.com/short-d/app/modern/mdrequest"
	"github.com/short-d/app/modern/mdrouting"
	"github.com/short-d/app/modern/mdruntime"
	"github.com/short-d/app/modern/mdservice"
	"github.com/short-d/app/modern/mdtimer"
	"github.com/short-d/app/modern/mdtracer"
	"github.com/short-d/short/app/adapter/db"
	"github.com/short-d/short/app/adapter/facebook"
	"github.com/short-d/short/app/adapter/github"
	"github.com/short-d/short/app/adapter/google"
	"github.com/short-d/short/app/adapter/graphql"
	"github.com/short-d/short/app/usecase/account"
	"github.com/short-d/short/app/usecase/requester"
	"github.com/short-d/short/app/usecase/url"
	"github.com/short-d/short/app/usecase/validator"
	"github.com/short-d/short/dep/provider"
)

// Injectors from wire.go:

func InjectCommandFactory() fw.CommandFactory {
	cobraFactory := mdcli.NewCobraFactory()
	return cobraFactory
}

func InjectDBConnector() fw.DBConnector {
	postgresConnector := mddb.NewPostgresConnector()
	return postgresConnector
}

func InjectDBMigrationTool() fw.DBMigrationTool {
	postgresMigrationTool := mddb.NewPostgresMigrationTool()
	return postgresMigrationTool
}

func InjectEnvironment() fw.Environment {
	goDotEnv := mdenv.NewGoDotEnv()
	return goDotEnv
}

func InjectGraphQLService(name string, prefix provider.LogPrefix, logLevel fw.LogLevel, sqlDB *sql.DB, graphqlPath provider.GraphQlPath, secret provider.ReCaptchaSecret, jwtSecret provider.JwtSecret, bufferSize provider.KeyGenBufferSize, kgsRPCConfig provider.KgsRPCConfig) (mdservice.Service, error) {
	stdOut := mdio.NewBuildInStdOut()
	timer := mdtimer.NewTimer()
	buildIn := mdruntime.NewBuildIn()
	local := provider.NewLocalLogger(prefix, logLevel, stdOut, timer, buildIn)
	tracer := mdtracer.NewLocal()
	urlSql := db.NewURLSql(sqlDB)
	retrieverPersist := url.NewRetrieverPersist(urlSql)
	userURLRelationSQL := db.NewUserURLRelationSQL(sqlDB)
	rpc, err := provider.NewKgsRPC(kgsRPCConfig)
	if err != nil {
		return mdservice.Service{}, err
	}
	remote, err := provider.NewRemote(bufferSize, rpc)
	if err != nil {
		return mdservice.Service{}, err
	}
	longLink := validator.NewLongLink()
	customAlias := validator.NewCustomAlias()
	creatorPersist := url.NewCreatorPersist(urlSql, userURLRelationSQL, remote, longLink, customAlias)
	client := mdhttp.NewClient()
	http := mdrequest.NewHTTP(client)
	reCaptcha := provider.NewReCaptchaService(http, secret)
	verifier := requester.NewVerifier(reCaptcha)
	cryptoTokenizer := provider.NewJwtGo(jwtSecret)
	tokenValidDuration := _wireTokenValidDurationValue
	authenticator := provider.NewAuthenticator(cryptoTokenizer, timer, tokenValidDuration)
	short := graphql.NewShort(local, tracer, retrieverPersist, creatorPersist, verifier, authenticator)
	server := provider.NewGraphGophers(graphqlPath, local, tracer, short)
	service := mdservice.New(name, server, local)
	return service, nil
}

var (
	_wireTokenValidDurationValue = provider.TokenValidDuration(oneDay)
)

func InjectRoutingService(name string, prefix provider.LogPrefix, logLevel fw.LogLevel, sqlDB *sql.DB, githubClientID provider.GithubClientID, githubClientSecret provider.GithubClientSecret, facebookClientID provider.FacebookClientID, facebookClientSecret provider.FacebookClientSecret, facebookRedirectURI provider.FacebookRedirectURI, googleClientID provider.GoogleClientID, googleClientSecret provider.GoogleClientSecret, googleRedirectURI provider.GoogleRedirectURI, jwtSecret provider.JwtSecret, webFrontendURL provider.WebFrontendURL) mdservice.Service {
	stdOut := mdio.NewBuildInStdOut()
	timer := mdtimer.NewTimer()
	buildIn := mdruntime.NewBuildIn()
	local := provider.NewLocalLogger(prefix, logLevel, stdOut, timer, buildIn)
	tracer := mdtracer.NewLocal()
	urlSql := db.NewURLSql(sqlDB)
	retrieverPersist := url.NewRetrieverPersist(urlSql)
	client := mdhttp.NewClient()
	http := mdrequest.NewHTTP(client)
	identityProvider := provider.NewGithubIdentityProvider(http, githubClientID, githubClientSecret)
	graphQL := mdrequest.NewGraphQL(http)
	githubAccount := github.NewAccount(graphQL)
	api := github.NewAPI(identityProvider, githubAccount)
	facebookIdentityProvider := provider.NewFacebookIdentityProvider(http, facebookClientID, facebookClientSecret, facebookRedirectURI)
	facebookAccount := facebook.NewAccount(http)
	facebookAPI := facebook.NewAPI(facebookIdentityProvider, facebookAccount)
	googleIdentityProvider := provider.NewGoogleIdentityProvider(http, googleClientID, googleClientSecret, googleRedirectURI)
	googleAccount := google.NewAccount(http)
	googleAPI := google.NewAPI(googleIdentityProvider, googleAccount)
	cryptoTokenizer := provider.NewJwtGo(jwtSecret)
	tokenValidDuration := _wireTokenValidDurationValue
	authenticator := provider.NewAuthenticator(cryptoTokenizer, timer, tokenValidDuration)
	userSQL := db.NewUserSQL(sqlDB)
	accountProvider := account.NewProvider(userSQL, timer)
	v := provider.NewShortRoutes(local, tracer, webFrontendURL, timer, retrieverPersist, api, facebookAPI, googleAPI, authenticator, accountProvider)
	server := mdrouting.NewBuiltIn(local, tracer, v)
	service := mdservice.New(name, server, local)
	return service
}

// wire.go:

const oneDay = 24 * time.Hour

var authSet = wire.NewSet(provider.NewJwtGo, wire.Value(provider.TokenValidDuration(oneDay)), provider.NewAuthenticator)

var observabilitySet = wire.NewSet(wire.Bind(new(fw.Logger), new(mdlogger.Local)), provider.NewLocalLogger, mdtracer.NewLocal)

var githubAPISet = wire.NewSet(provider.NewGithubIdentityProvider, github.NewAccount, github.NewAPI)

var facebookAPISet = wire.NewSet(provider.NewFacebookIdentityProvider, facebook.NewAccount, facebook.NewAPI)

var googleAPISet = wire.NewSet(provider.NewGoogleIdentityProvider, google.NewAccount, google.NewAPI)
