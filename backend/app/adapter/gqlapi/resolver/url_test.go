// +build !integration all

package resolver

import (
	"testing"
	"time"

	"github.com/short-d/app/fw/assert"
	"github.com/short-d/short/backend/app/adapter/gqlapi/scalar"
	"github.com/short-d/short/backend/app/entity"
)

func TestURL_Alias(t *testing.T) {
	t.Parallel()
	urlResolver := URL{url: entity.URL{Alias: "TestAlias"}}

	expected := urlResolver.url.Alias
	got := *urlResolver.Alias()
	assert.Equal(t, got, expected, "*urlTest.Alias() = %v; want %v", expected, got)
}

func TestURL_OriginalURL(t *testing.T) {
	t.Parallel()
	urlResolver := URL{url: entity.URL{OriginalURL: "TestOriginalUrl"}}

	expected := urlResolver.url.OriginalURL
	got := *urlResolver.OriginalURL()
	assert.Equal(t, got, expected, "*urlResolver.OriginalURL() = %v; want %v", expected, got)
}

func TestURL_ExpireAt(t *testing.T) {
	t.Parallel()
	timeAfter := time.Now().Add(5 * time.Second)
	testCases := []struct {
		url      URL
		expected *scalar.Time
	}{
		{
			url:      URL{url: entity.URL{ExpireAt: &timeAfter}},
			expected: &scalar.Time{Time: timeAfter},
		},
		{
			url: URL{url: entity.URL{ExpireAt: nil}},
		},
	}
	for _, testCase := range testCases {
		testCase := testCase
		assert.Equal(t, testCase.expected, testCase.url.ExpireAt())
	}
}
