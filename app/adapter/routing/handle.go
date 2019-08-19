package routing

import (
	"fmt"
	"net/http"
	"short/app/adapter/account"
	"short/app/adapter/oauth"
	"short/app/usecase/auth"
	"short/app/usecase/url"
	"short/fw"
	"strings"
)

func NewOriginalUrl(
	logger fw.Logger,
	tracer fw.Tracer,
	urlRetriever url.Retriever,
	timer fw.Timer,
) fw.Handle {
	return func(w http.ResponseWriter, r *http.Request, params fw.Params) {
		trace := tracer.BeginTrace("OriginalUrl")

		alias := params["alias"]

		trace1 := trace.Next("GetUrlAfter")
		u, err := urlRetriever.GetAfter(trace1, alias, timer.Now())
		trace1.End()

		if err != nil {
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			logger.Error(err)
			return
		}

		originUrl := u.OriginalUrl
		http.Redirect(w, r, originUrl, http.StatusSeeOther)
		trace.End()
	}
}

func getFilenameFromPath(path string, indexFile string) string {
	filePath := strings.Trim(path, "/")
	if filePath == "" {
		return indexFile
	}
	return filePath
}

func NewServeFile(logger fw.Logger, tracer fw.Tracer, wwwRoot string) fw.Handle {
	fs := http.FileServer(http.Dir(wwwRoot))

	return func(w http.ResponseWriter, r *http.Request, params fw.Params) {
		fileName := getFilenameFromPath(r.URL.Path, "index.html")
		logger.Info(fmt.Sprintf("serving %s from %s", fileName, wwwRoot))

		fs.ServeHTTP(w, r)
	}
}

func NewGithubSignIn(
	logger fw.Logger,
	tracer fw.Tracer,
	githubOAuth oauth.Github,
	authenticator auth.Authenticator,
) fw.Handle {
	return func(w http.ResponseWriter, r *http.Request, params fw.Params) {
		token := getToken(r, params)
		if authenticator.IsSignedIn(token) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		signInLink := githubOAuth.GetAuthorizationUrl()
		http.Redirect(w, r, signInLink, http.StatusSeeOther)
	}
}

func NewGithubSignInCallback(
	logger fw.Logger,
	tracer fw.Tracer,
	githubOAuth oauth.Github,
	githubAccount account.Github,
	authenticator auth.Authenticator,
) fw.Handle {
	return func(w http.ResponseWriter, r *http.Request, params fw.Params) {
		code := params["code"]
		if len(code) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		accessToken, _, err := githubOAuth.RequestAccessToken(code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		email, err := githubAccount.GetEmail(accessToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		authToken, err := authenticator.GenerateToken(email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w = setToken(w, authToken)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
