package app

import (
	"time"

	"github.com/short-d/app/fw"
	"github.com/short-d/short/dep"
	"github.com/short-d/short/dep/provider"
)

// ServiceConfig represents require parameters for the backend APIs
type ServiceConfig struct {
	ServerEnv            string
	LogPrefix            string
	LogLevel             fw.LogLevel
	MigrationRoot        string
	RecaptchaSecret      string
	GithubClientID       string
	GithubClientSecret   string
	FacebookClientID     string
	FacebookClientSecret string
	FacebookRedirectURI  string
	GoogleClientID       string
	GoogleClientSecret   string
	GoogleRedirectURI    string
	JwtSecret            string
	WebFrontendURL       string
	GraphQLAPIPort       int
	HTTPAPIPort          int
	KeyGenBufferSize     int
	KgsHostname          string
	KgsPort              int
	AuthTokenLifetime    time.Duration
	DataDogAPIKey        string
	SegmentAPIKey        string
	IPStackAPIKey        string
	GoogleAPIKey         string
}

// Start launches the GraphQL & HTTP APIs
func Start(
	dbConfig fw.DBConfig,
	config ServiceConfig,
	dbConnector fw.DBConnector,
	dbMigrationTool fw.DBMigrationTool,
) {
	db, err := dbConnector.Connect(dbConfig)
	if err != nil {
		panic(err)
	}

	err = dbMigrationTool.MigrateUp(db, config.MigrationRoot)
	if err != nil {
		panic(err)
	}

	serverEnv := fw.ServerEnv(config.ServerEnv)
	kgsBufferSize := provider.KeyGenBufferSize(config.KeyGenBufferSize)
	kgsRPCConfig := provider.KgsRPCConfig{
		Hostname: config.KgsHostname,
		Port:     config.KgsPort,
	}

	dataDogAPIKey := provider.DataDogAPIKey(config.DataDogAPIKey)
	segmentAPIKey := provider.SegmentAPIKey(config.SegmentAPIKey)
	ipStackAPIKey := provider.IPStackAPIKey(config.IPStackAPIKey)
	googleAPIKey := provider.GoogleAPIKey(config.GoogleAPIKey)

	graphqlAPI, err := dep.InjectGraphQLService(
		"GraphQL API",
		serverEnv,
		provider.LogPrefix(config.LogPrefix),
		config.LogLevel,
		db,
		"/graphql",
		provider.ReCaptchaSecret(config.RecaptchaSecret),
		provider.JwtSecret(config.JwtSecret),
		kgsBufferSize,
		kgsRPCConfig,
		provider.TokenValidDuration(config.AuthTokenLifetime),
		dataDogAPIKey,
		segmentAPIKey,
		ipStackAPIKey,
		googleAPIKey,
	)
	if err != nil {
		panic(err)
	}

	graphqlAPI.Start(config.GraphQLAPIPort)

	httpAPI, err := dep.InjectRoutingService(
		"Routing API",
		serverEnv,
		provider.LogPrefix(config.LogPrefix),
		config.LogLevel,
		db,
		provider.GithubClientID(config.GithubClientID),
		provider.GithubClientSecret(config.GithubClientSecret),
		provider.FacebookClientID(config.FacebookClientID),
		provider.FacebookClientSecret(config.FacebookClientSecret),
		provider.FacebookRedirectURI(config.FacebookRedirectURI),
		provider.GoogleClientID(config.GoogleClientID),
		provider.GoogleClientSecret(config.GoogleClientSecret),
		provider.GoogleRedirectURI(config.GoogleRedirectURI),
		provider.JwtSecret(config.JwtSecret),
		kgsBufferSize,
		kgsRPCConfig,
		provider.WebFrontendURL(config.WebFrontendURL),
		provider.TokenValidDuration(config.AuthTokenLifetime),
		dataDogAPIKey,
		segmentAPIKey,
		ipStackAPIKey,
	)
	if err != nil {
		panic(err)
	}

	httpAPI.StartAndWait(config.HTTPAPIPort)
}
