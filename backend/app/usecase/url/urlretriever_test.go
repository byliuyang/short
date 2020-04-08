// +build !integration all

package url

import (
	"testing"
	"time"

	"github.com/short-d/app/mdtest"
	"github.com/short-d/short/app/entity"
	"github.com/short-d/short/app/usecase/repository"
)

type urlMap = map[string]entity.URL

func TestUrlRetriever_GetURL(t *testing.T) {
	t.Parallel()

	now := time.Now()
	before := now.Add(-5 * time.Second)
	after := now.Add(5 * time.Second)

	testCases := []struct {
		name        string
		urls        urlMap
		alias       string
		expiringAt  *time.Time
		hasErr      bool
		expectedURL entity.URL
	}{
		{
			name:        "alias not found",
			urls:        urlMap{},
			alias:       "220uFicCJj",
			expiringAt:  &now,
			hasErr:      true,
			expectedURL: entity.URL{},
		},
		{
			name: "url expired",
			urls: urlMap{
				"220uFicCJj": entity.URL{
					Alias:    "220uFicCJj",
					ExpireAt: &before,
				},
			},
			alias:       "220uFicCJj",
			expiringAt:  &now,
			hasErr:      true,
			expectedURL: entity.URL{},
		},
		{
			name: "url never expire",
			urls: urlMap{
				"220uFicCJj": entity.URL{
					Alias:    "220uFicCJj",
					ExpireAt: nil,
				},
			},
			alias:      "220uFicCJj",
			expiringAt: &now,
			hasErr:     false,
			expectedURL: entity.URL{
				Alias:    "220uFicCJj",
				ExpireAt: nil,
			},
		},
		{
			name: "unexpired url found",
			urls: urlMap{
				"220uFicCJj": entity.URL{
					Alias:    "220uFicCJj",
					ExpireAt: &after,
				},
			},
			alias:      "220uFicCJj",
			expiringAt: &now,
			hasErr:     false,
			expectedURL: entity.URL{
				Alias:    "220uFicCJj",
				ExpireAt: &after,
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			fakeURLRepo := repository.NewURLFake(testCase.urls)
			fakeUserURLRelationRepo := repository.NewUserURLRepoFake([]entity.User{}, []entity.URL{})
			retriever := NewRetrieverPersist(&fakeURLRepo, &fakeUserURLRelationRepo)
			url, err := retriever.GetURL(testCase.alias, testCase.expiringAt)

			if testCase.hasErr {
				mdtest.NotEqual(t, nil, err)
				return
			}
			mdtest.Equal(t, nil, err)
			mdtest.Equal(t, testCase.expectedURL, url)
		})
	}
}

func TestRetrieverPersist_GetURLs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		urls         urlMap
		users        []entity.User
		createdURLs  []entity.URL
		user         entity.User
		hasErr       bool
		expectedURLs []entity.URL
	}{
		{
			name: "user created URLs",
			urls: urlMap{
				"google": entity.URL{
					Alias:       "google",
					OriginalURL: "https://www.google.com/",
				},
				"short": entity.URL{
					Alias:       "short",
					OriginalURL: "https://github.com/short-d/short/",
				},
				"mozilla": entity.URL{
					Alias:       "mozilla",
					OriginalURL: "https://www.mozilla.org/",
				},
			},
			users: []entity.User{
				{
					ID:    "12345",
					Name:  "Test User",
					Email: "test@gmail.com",
				}, {
					ID:    "12345",
					Name:  "Test User",
					Email: "test@gmail.com",
				}, {
					ID:    "12346",
					Name:  "Test User 2",
					Email: "test2@gmail.com",
				},
			},
			createdURLs: []entity.URL{
				{
					Alias:       "google",
					OriginalURL: "https://www.google.com/",
				},
				{
					Alias:       "short",
					OriginalURL: "https://github.com/short-d/short/",
				},
				{
					Alias:       "mozilla",
					OriginalURL: "https://www.mozilla.org/",
				},
			},
			user: entity.User{
				ID:    "12345",
				Name:  "Test User",
				Email: "test@gmail.com",
			},
			hasErr: false,
			expectedURLs: []entity.URL{
				{
					Alias:       "google",
					OriginalURL: "https://www.google.com/",
				},
				{
					Alias:       "short",
					OriginalURL: "https://github.com/short-d/short/",
				},
			},
		},
		{
			name: "user has no URL",
			urls: urlMap{
				"google": entity.URL{
					Alias:       "google",
					OriginalURL: "https://www.google.com/",
				},
				"short": entity.URL{
					Alias:       "short",
					OriginalURL: "https://github.com/short-d/short/",
				},
				"mozilla": entity.URL{
					Alias:       "mozilla",
					OriginalURL: "https://www.mozilla.org/",
				},
			},
			users: []entity.User{
				{
					ID:    "12345",
					Name:  "Test User",
					Email: "test@gmail.com",
				}, {
					ID:    "12345",
					Name:  "Test User",
					Email: "test@gmail.com",
				}, {
					ID:    "12345",
					Name:  "Test User",
					Email: "test@gmail.com",
				},
			},
			createdURLs: []entity.URL{
				{
					Alias:       "google",
					OriginalURL: "https://www.google.com/",
				},
				{
					Alias:       "short",
					OriginalURL: "https://github.com/short-d/short/",
				},
				{
					Alias:       "mozilla",
					OriginalURL: "https://www.mozilla.org/",
				},
			},
			user: entity.User{
				ID:    "12346",
				Name:  "Test User 2",
				Email: "test2@gmail.com",
			},
			hasErr:       false,
			expectedURLs: []entity.URL{},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			fakeURLRepo := repository.NewURLFake(testCase.urls)
			fakeUserURLRelationRepo := repository.NewUserURLRepoFake(testCase.users, testCase.createdURLs)
			retriever := NewRetrieverPersist(&fakeURLRepo, &fakeUserURLRelationRepo)

			urls, err := retriever.GetURLsByUser(testCase.user)
			if testCase.hasErr {
				mdtest.NotEqual(t, nil, err)
				return
			}

			mdtest.Equal(t, nil, err)
			mdtest.Equal(t, testCase.expectedURLs, urls)
		})
	}
}
