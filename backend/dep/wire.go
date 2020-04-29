//+build wireinject

package dep

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/short-d/app/fw"
	"github.com/short-d/app/modern/mdanalytics"
	"github.com/short-d/app/modern/mdcli"
	"github.com/short-d/app/modern/mddb"
	"github.com/short-d/app/modern/mdenv"
	"github.com/short-d/app/modern/mdgeo"
	"github.com/short-d/app/modern/mdhttp"
	"github.com/short-d/app/modern/mdio"
	"github.com/short-d/app/modern/mdlogger"
	"github.com/short-d/app/modern/mdmetrics"
	"github.com/short-d/app/modern/mdnetwork"
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
	"github.com/short-d/short/app/adapter/kgs"
	"github.com/short-d/short/app/adapter/request"
	"github.com/short-d/short/app/usecase/account"
	"github.com/short-d/short/app/usecase/changelog"
	"github.com/short-d/short/app/usecase/repository"
	"github.com/short-d/short/app/usecase/requester"
	"github.com/short-d/short/app/usecase/risk"
	"github.com/short-d/short/app/usecase/service"
	"github.com/short-d/short/app/usecase/url"
	"github.com/short-d/short/app/usecase/validator"
	"github.com/short-d/short/dep/provider"
)

var authSet = wire.NewSet(
	provider.NewJwtGo,
	provider.NewAuthenticator,
)

var observabilitySet = wire.NewSet(
	wire.Bind(new(fw.StdOut), new(mdio.StdOut)),
	wire.Bind(new(fw.Logger), new(mdlogger.Logger)),
	wire.Bind(new(fw.Metrics), new(mdmetrics.DataDog)),
	wire.Bind(new(fw.Analytics), new(mdanalytics.Segment)),
	wire.Bind(new(fw.Network), new(mdnetwork.Proxy)),

	mdio.NewBuildInStdOut,
	provider.NewEntryRepositorySwitch,
	provider.NewLogger,
	mdtracer.NewLocal,
	provider.NewDataDogMetrics,
	provider.NewSegment,
	mdnetwork.NewProxy,
	request.NewClient,
	request.NewInstrumentationFactory,
)

var githubAPISet = wire.NewSet(
	provider.NewGithubIdentityProvider,
	github.NewAccount,
	github.NewAPI,
)

var facebookAPISet = wire.NewSet(
	provider.NewFacebookIdentityProvider,
	facebook.NewAccount,
	facebook.NewAPI,
)

var googleAPISet = wire.NewSet(
	provider.NewGoogleIdentityProvider,
	google.NewAccount,
	google.NewAPI,
)

var keyGenSet = wire.NewSet(
	wire.Bind(new(service.KeyFetcher), new(kgs.RPC)),
	provider.NewKgsRPC,
	provider.NewKeyGenerator,
)

var featureDecisionSet = wire.NewSet(
	wire.Bind(new(repository.FeatureToggle), new(db.FeatureToggleSQL)),
	db.NewFeatureToggleSQL,
	provider.NewFeatureDecisionMakerFactorySwitch,
)

// InjectCommandFactory creates CommandFactory with configured dependencies.
func InjectCommandFactory() fw.CommandFactory {
	wire.Build(
		wire.Bind(new(fw.CommandFactory), new(mdcli.CobraFactory)),
		mdcli.NewCobraFactory,
	)
	return mdcli.CobraFactory{}
}

// InjectDBConnector creates DBConnector with configured dependencies.
func InjectDBConnector() fw.DBConnector {
	wire.Build(
		wire.Bind(new(fw.DBConnector), new(mddb.PostgresConnector)),
		mddb.NewPostgresConnector,
	)
	return mddb.PostgresConnector{}
}

// InjectDBMigrationTool creates DBMigrationTool with configured dependencies.
func InjectDBMigrationTool() fw.DBMigrationTool {
	wire.Build(
		wire.Bind(new(fw.DBMigrationTool), new(mddb.PostgresMigrationTool)),
		mddb.NewPostgresMigrationTool,
	)
	return mddb.PostgresMigrationTool{}
}

// InjectEnvironment creates Environment with configured dependencies.
func InjectEnvironment() fw.Environment {
	wire.Build(
		wire.Bind(new(fw.Environment), new(mdenv.GoDotEnv)),
		mdenv.NewGoDotEnv,
	)
	return mdenv.GoDotEnv{}
}

// InjectGraphQLService creates GraphQL service with configured dependencies.
func InjectGraphQLService(
	name string,
	serverEnv fw.ServerEnv,
	prefix provider.LogPrefix,
	logLevel fw.LogLevel,
	sqlDB *sql.DB,
	graphqlPath provider.GraphQlPath,
	secret provider.ReCaptchaSecret,
	jwtSecret provider.JwtSecret,
	bufferSize provider.KeyGenBufferSize,
	kgsRPCConfig provider.KgsRPCConfig,
	tokenValidDuration provider.TokenValidDuration,
	dataDogAPIKey provider.DataDogAPIKey,
	segmentAPIKey provider.SegmentAPIKey,
	ipStackAPIKey provider.IPStackAPIKey,
	googleAPIKey provider.GoogleAPIKey,
) (mdservice.Service, error) {
	wire.Build(
		wire.Bind(new(fw.ProgramRuntime), new(mdruntime.BuildIn)),
		wire.Bind(new(fw.GraphQLAPI), new(graphql.Short)),
		wire.Bind(new(changelog.ChangeLog), new(changelog.Persist)),
		wire.Bind(new(url.Retriever), new(url.RetrieverPersist)),
		wire.Bind(new(url.Creator), new(url.CreatorPersist)),
		wire.Bind(new(repository.UserURLRelation), new(db.UserURLRelationSQL)),
		wire.Bind(new(repository.ChangeLog), new(db.ChangeLogSQL)),
		wire.Bind(new(repository.URL), new(*db.URLSql)),
		wire.Bind(new(risk.BlackList), new(google.SafeBrowsing)),
		wire.Bind(new(fw.HTTPRequest), new(mdrequest.HTTP)),

		observabilitySet,
		authSet,
		keyGenSet,

		mdruntime.NewBuildIn,
		mdservice.New,
		provider.NewGraphGophers,
		mdhttp.NewClient,
		mdrequest.NewHTTP,
		mdtimer.NewTimer,
		provider.NewSafeBrowsing,
		risk.NewDetector,

		db.NewChangeLogSQL,
		db.NewURLSql,
		db.NewUserURLRelationSQL,
		validator.NewLongLink,
		validator.NewCustomAlias,
		changelog.NewPersist,
		url.NewRetrieverPersist,
		url.NewCreatorPersist,
		provider.NewReCaptchaService,
		requester.NewVerifier,
		graphql.NewShort,
	)
	return mdservice.Service{}, nil
}

// InjectRoutingService creates routing service with configured dependencies.
func InjectRoutingService(
	name string,
	serverEnv fw.ServerEnv,
	prefix provider.LogPrefix,
	logLevel fw.LogLevel,
	sqlDB *sql.DB,
	githubClientID provider.GithubClientID,
	githubClientSecret provider.GithubClientSecret,
	facebookClientID provider.FacebookClientID,
	facebookClientSecret provider.FacebookClientSecret,
	facebookRedirectURI provider.FacebookRedirectURI,
	googleClientID provider.GoogleClientID,
	googleClientSecret provider.GoogleClientSecret,
	googleRedirectURI provider.GoogleRedirectURI,
	jwtSecret provider.JwtSecret,
	bufferSize provider.KeyGenBufferSize,
	kgsRPCConfig provider.KgsRPCConfig,
	webFrontendURL provider.WebFrontendURL,
	tokenValidDuration provider.TokenValidDuration,
	dataDogAPIKey provider.DataDogAPIKey,
	segmentAPIKey provider.SegmentAPIKey,
	ipStackAPIKey provider.IPStackAPIKey,
) (mdservice.Service, error) {
	wire.Build(
		wire.Bind(new(fw.ProgramRuntime), new(mdruntime.BuildIn)),
		wire.Bind(new(url.Retriever), new(url.RetrieverPersist)),
		wire.Bind(new(repository.UserURLRelation), new(db.UserURLRelationSQL)),
		wire.Bind(new(repository.User), new(*db.UserSQL)),
		wire.Bind(new(repository.URL), new(*db.URLSql)),
		wire.Bind(new(fw.HTTPRequest), new(mdrequest.HTTP)),
		wire.Bind(new(fw.GraphQlRequest), new(mdrequest.GraphQL)),
		wire.Bind(new(fw.GeoLocation), new(mdgeo.IPStack)),

		observabilitySet,
		authSet,
		githubAPISet,
		facebookAPISet,
		googleAPISet,
		keyGenSet,
		featureDecisionSet,

		mdruntime.NewBuildIn,
		mdservice.New,
		mdrouting.NewBuiltIn,
		mdhttp.NewClient,
		mdrequest.NewHTTP,
		mdrequest.NewGraphQL,
		mdtimer.NewTimer,
		provider.NewIPStack,

		db.NewUserSQL,
		db.NewURLSql,
		db.NewUserURLRelationSQL,
		url.NewRetrieverPersist,
		account.NewProvider,
		provider.NewShortRoutes,
	)
	return mdservice.Service{}, nil
}
