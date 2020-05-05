// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package dep

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/short-d/app/fw/analytics"
	"github.com/short-d/app/fw/cli"
	"github.com/short-d/app/fw/db"
	"github.com/short-d/app/fw/env"
	"github.com/short-d/app/fw/graphql"
	"github.com/short-d/app/fw/io"
	"github.com/short-d/app/fw/logger"
	"github.com/short-d/app/fw/metrics"
	"github.com/short-d/app/fw/network"
	"github.com/short-d/app/fw/runtime"
	"github.com/short-d/app/fw/service"
	"github.com/short-d/app/fw/timer"
	"github.com/short-d/app/fw/webreq"
	"github.com/short-d/short/backend/app/adapter/facebook"
	"github.com/short-d/short/backend/app/adapter/github"
	"github.com/short-d/short/backend/app/adapter/google"
	"github.com/short-d/short/backend/app/adapter/gqlapi"
	"github.com/short-d/short/backend/app/adapter/kgs"
	"github.com/short-d/short/backend/app/adapter/request"
	"github.com/short-d/short/backend/app/adapter/sqldb"
	"github.com/short-d/short/backend/app/usecase/account"
	"github.com/short-d/short/backend/app/usecase/changelog"
	"github.com/short-d/short/backend/app/usecase/external"
	"github.com/short-d/short/backend/app/usecase/repository"
	"github.com/short-d/short/backend/app/usecase/requester"
	"github.com/short-d/short/backend/app/usecase/risk"
	"github.com/short-d/short/backend/app/usecase/url"
	"github.com/short-d/short/backend/app/usecase/validator"
	"github.com/short-d/short/backend/dep/provider"
)

// Injectors from wire.go:

func InjectCommandFactory() cli.CommandFactory {
	cobraFactory := cli.NewCobraFactory()
	return cobraFactory
}

func InjectDBConnector() db.Connector {
	postgresConnector := db.NewPostgresConnector()
	return postgresConnector
}

func InjectDBMigrationTool() db.MigrationTool {
	postgresMigrationTool := db.NewPostgresMigrationTool()
	return postgresMigrationTool
}

func InjectEnv() env.Env {
	goDotEnv := env.NewGoDotEnv()
	return goDotEnv
}

func InjectGraphQLService(runtime2 env.Runtime, prefix provider.LogPrefix, logLevel logger.LogLevel, sqlDB *sql.DB, graphqlPath provider.GraphQLPath, secret provider.ReCaptchaSecret, jwtSecret provider.JwtSecret, bufferSize provider.KeyGenBufferSize, kgsRPCConfig provider.KgsRPCConfig, tokenValidDuration provider.TokenValidDuration, dataDogAPIKey provider.DataDogAPIKey, segmentAPIKey provider.SegmentAPIKey, ipStackAPIKey provider.IPStackAPIKey, googleAPIKey provider.GoogleAPIKey) (service.GraphQL, error) {
	system := timer.NewSystem()
	program := runtime.NewProgram()
	deployment := env.NewDeployment(runtime2)
	stdOut := io.NewStdOut()
	client := webreq.NewHTTPClient()
	http := webreq.NewHTTP(client)
	entryRepository := provider.NewEntryRepositorySwitch(runtime2, deployment, stdOut, dataDogAPIKey, http)
	loggerLogger := provider.NewLogger(prefix, logLevel, system, program, entryRepository)
	urlSql := sqldb.NewURLSql(sqlDB)
	userURLRelationSQL := sqldb.NewUserURLRelationSQL(sqlDB)
	retrieverPersist := url.NewRetrieverPersist(urlSql, userURLRelationSQL)
	rpc, err := provider.NewKgsRPC(kgsRPCConfig)
	if err != nil {
		return service.GraphQL{}, err
	}
	keyGenerator, err := provider.NewKeyGenerator(bufferSize, rpc)
	if err != nil {
		return service.GraphQL{}, err
	}
	longLink := validator.NewLongLink()
	customAlias := validator.NewCustomAlias()
	safeBrowsing := provider.NewSafeBrowsing(googleAPIKey, http)
	detector := risk.NewDetector(safeBrowsing)
	creatorPersist := url.NewCreatorPersist(urlSql, userURLRelationSQL, keyGenerator, longLink, customAlias, system, detector)
	changeLogSQL := sqldb.NewChangeLogSQL(sqlDB)
	persist := changelog.NewPersist(keyGenerator, system, changeLogSQL)
	reCaptcha := provider.NewReCaptchaService(http, secret)
	verifier := requester.NewVerifier(reCaptcha)
	tokenizer := provider.NewJwtGo(jwtSecret)
	authenticator := provider.NewAuthenticator(tokenizer, system, tokenValidDuration)
	short := gqlapi.NewShort(loggerLogger, retrieverPersist, creatorPersist, persist, verifier, authenticator)
	graphGopherHandler := graphql.NewGraphGopherHandler(short)
	graphQL := provider.NewGraphQLService(graphqlPath, graphGopherHandler, loggerLogger)
	return graphQL, nil
}

func InjectRoutingService(runtime2 env.Runtime, prefix provider.LogPrefix, logLevel logger.LogLevel, sqlDB *sql.DB, githubClientID provider.GithubClientID, githubClientSecret provider.GithubClientSecret, facebookClientID provider.FacebookClientID, facebookClientSecret provider.FacebookClientSecret, facebookRedirectURI provider.FacebookRedirectURI, googleClientID provider.GoogleClientID, googleClientSecret provider.GoogleClientSecret, googleRedirectURI provider.GoogleRedirectURI, jwtSecret provider.JwtSecret, bufferSize provider.KeyGenBufferSize, kgsRPCConfig provider.KgsRPCConfig, webFrontendURL provider.WebFrontendURL, tokenValidDuration provider.TokenValidDuration, dataDogAPIKey provider.DataDogAPIKey, segmentAPIKey provider.SegmentAPIKey, ipStackAPIKey provider.IPStackAPIKey) (service.Routing, error) {
	system := timer.NewSystem()
	program := runtime.NewProgram()
	deployment := env.NewDeployment(runtime2)
	stdOut := io.NewStdOut()
	client := webreq.NewHTTPClient()
	http := webreq.NewHTTP(client)
	entryRepository := provider.NewEntryRepositorySwitch(runtime2, deployment, stdOut, dataDogAPIKey, http)
	loggerLogger := provider.NewLogger(prefix, logLevel, system, program, entryRepository)
	dataDog := provider.NewDataDogMetrics(dataDogAPIKey, http, system, runtime2)
	segment := provider.NewSegment(segmentAPIKey, system, loggerLogger)
	rpc, err := provider.NewKgsRPC(kgsRPCConfig)
	if err != nil {
		return service.Routing{}, err
	}
	keyGenerator, err := provider.NewKeyGenerator(bufferSize, rpc)
	if err != nil {
		return service.Routing{}, err
	}
	proxy := network.NewProxy()
	ipStack := provider.NewIPStack(ipStackAPIKey, http, loggerLogger)
	requestClient := request.NewClient(proxy, ipStack)
	instrumentationFactory := request.NewInstrumentationFactory(loggerLogger, system, dataDog, segment, keyGenerator, requestClient)
	urlSql := sqldb.NewURLSql(sqlDB)
	userURLRelationSQL := sqldb.NewUserURLRelationSQL(sqlDB)
	retrieverPersist := url.NewRetrieverPersist(urlSql, userURLRelationSQL)
	identityProvider := provider.NewGithubIdentityProvider(http, githubClientID, githubClientSecret)
	clientFactory := graphql.NewClientFactory(http)
	githubAccount := github.NewAccount(clientFactory)
	api := github.NewAPI(identityProvider, githubAccount)
	facebookIdentityProvider := provider.NewFacebookIdentityProvider(http, facebookClientID, facebookClientSecret, facebookRedirectURI)
	facebookAccount := facebook.NewAccount(http)
	facebookAPI := facebook.NewAPI(facebookIdentityProvider, facebookAccount)
	googleIdentityProvider := provider.NewGoogleIdentityProvider(http, googleClientID, googleClientSecret, googleRedirectURI)
	googleAccount := google.NewAccount(http)
	googleAPI := google.NewAPI(googleIdentityProvider, googleAccount)
	featureToggleSQL := sqldb.NewFeatureToggleSQL(sqlDB)
	decisionMakerFactory := provider.NewFeatureDecisionMakerFactorySwitch(deployment, featureToggleSQL)
	tokenizer := provider.NewJwtGo(jwtSecret)
	authenticator := provider.NewAuthenticator(tokenizer, system, tokenValidDuration)
	userSQL := sqldb.NewUserSQL(sqlDB)
	accountProvider := account.NewProvider(userSQL, system)
	v := provider.NewShortRoutes(instrumentationFactory, webFrontendURL, system, retrieverPersist, api, facebookAPI, googleAPI, decisionMakerFactory, authenticator, accountProvider)
	routing := service.NewRouting(loggerLogger, v)
	return routing, nil
}

// wire.go:

var authSet = wire.NewSet(provider.NewJwtGo, provider.NewAuthenticator)

var observabilitySet = wire.NewSet(wire.Bind(new(io.Output), new(io.StdOut)), wire.Bind(new(runtime.Runtime), new(runtime.Program)), wire.Bind(new(metrics.Metrics), new(metrics.DataDog)), wire.Bind(new(analytics.Analytics), new(analytics.Segment)), wire.Bind(new(network.Network), new(network.Proxy)), io.NewStdOut, provider.NewEntryRepositorySwitch, provider.NewLogger, runtime.NewProgram, provider.NewDataDogMetrics, provider.NewSegment, network.NewProxy, request.NewClient, request.NewInstrumentationFactory)

var githubAPISet = wire.NewSet(provider.NewGithubIdentityProvider, github.NewAccount, github.NewAPI)

var facebookAPISet = wire.NewSet(provider.NewFacebookIdentityProvider, facebook.NewAccount, facebook.NewAPI)

var googleAPISet = wire.NewSet(provider.NewGoogleIdentityProvider, google.NewAccount, google.NewAPI)

var keyGenSet = wire.NewSet(wire.Bind(new(external.KeyFetcher), new(kgs.RPC)), provider.NewKgsRPC, provider.NewKeyGenerator)

var featureDecisionSet = wire.NewSet(wire.Bind(new(repository.FeatureToggle), new(sqldb.FeatureToggleSQL)), sqldb.NewFeatureToggleSQL, provider.NewFeatureDecisionMakerFactorySwitch)
