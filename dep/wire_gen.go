// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package dep

import (
	"database/sql"
	"short/app/adapter/account"
	"short/app/adapter/graphql"
	"short/app/adapter/repo"
	"short/app/usecase/keygen"
	"short/app/usecase/requester"
	"short/app/usecase/url"
	"short/dep/inject"
	"short/modern/mdhttp"
	"short/modern/mdlogger"
	"short/modern/mdrequest"
	"short/modern/mdrouting"
	"short/modern/mdservice"
	"short/modern/mdtimer"
	"short/modern/mdtracer"
)

// Injectors from wire.go:

func InitGraphQlService(name string, db *sql.DB, graphqlPath inject.GraphQlPath, secret inject.ReCaptchaSecret) mdservice.Service {
	logger := mdlogger.NewLocal()
	tracer := mdtracer.NewLocal()
	repoUrl := repo.NewUrlSql(db)
	retriever := url.NewRetrieverPersist(repoUrl)
	keyGenerator := keygen.NewInMemory()
	creator := url.NewCreatorPersist(repoUrl, keyGenerator)
	client := mdhttp.NewClient()
	httpRequest := mdrequest.NewHttp(client)
	reCaptcha := inject.ReCaptchaService(httpRequest, secret)
	verifier := requester.NewVerifier(reCaptcha)
	graphQlApi := graphql.NewShort(logger, tracer, retriever, creator, verifier)
	server := inject.GraphGophers(graphqlPath, logger, tracer, graphQlApi)
	service := mdservice.New(name, server, logger)
	return service
}

func InitRoutingService(name string, db *sql.DB, wwwRoot inject.WwwRoot, githubClientId inject.GithubClientId, githubClientSecret inject.GithubClientSecret, jwtSecret inject.JwtSecret) mdservice.Service {
	logger := mdlogger.NewLocal()
	tracer := mdtracer.NewLocal()
	timer := mdtimer.NewTimer()
	repoUrl := repo.NewUrlSql(db)
	retriever := url.NewRetrieverPersist(repoUrl)
	client := mdhttp.NewClient()
	httpRequest := mdrequest.NewHttp(client)
	github := inject.GithubOAuth(httpRequest, githubClientId, githubClientSecret)
	graphQlRequest := mdrequest.NewGraphQl(httpRequest)
	accountGithub := account.NewGithub(graphQlRequest)
	cryptoTokenizer := inject.JwtGo(jwtSecret)
	authenticator := inject.Authenticator(cryptoTokenizer, timer)
	v := inject.ShortRoutes(logger, tracer, wwwRoot, timer, retriever, github, accountGithub, authenticator)
	server := mdrouting.NewBuiltIn(logger, tracer, v)
	service := mdservice.New(name, server, logger)
	return service
}
