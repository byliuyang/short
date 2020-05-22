// +build integration all

package sqldb_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/short-d/app/fw/assert"
	"github.com/short-d/app/fw/db/dbtest"
	"github.com/short-d/short/backend/app/adapter/sqldb"
	"github.com/short-d/short/backend/app/adapter/sqldb/table"
	"github.com/short-d/short/backend/app/entity"
)

var insertShortLinkRowSQL = fmt.Sprintf(`
INSERT INTO %s (%s, %s, %s, %s, %s)
VALUES ($1, $2, $3, $4, $5)`,
	table.ShortLink.TableName,
	table.ShortLink.ColumnAlias,
	table.ShortLink.ColumnLongLink,
	table.ShortLink.ColumnCreatedAt,
	table.ShortLink.ColumnExpireAt,
	table.ShortLink.ColumnUpdatedAt,
)

type shortLinkTableRow struct {
	alias     string
	longLink  string
	createdAt *time.Time
	expireAt  *time.Time
	updatedAt *time.Time
}

func TestShortLinkSql_IsAliasExist(t *testing.T) {
	testCases := []struct {
		name       string
		tableRows  []shortLinkTableRow
		alias      string
		expIsExist bool
	}{
		{
			name:       "alias doesn't exist",
			alias:      "gg",
			tableRows:  []shortLinkTableRow{},
			expIsExist: false,
		},
		{
			name:  "alias found",
			alias: "gg",
			tableRows: []shortLinkTableRow{
				{alias: "gg"},
			},
			expIsExist: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dbtest.AccessTestDB(
				dbConnector,
				dbMigrationTool,
				dbMigrationRoot,
				dbConfig,
				func(sqlDB *sql.DB) {
					insertShortLinkTableRows(t, sqlDB, testCase.tableRows)

					shortLinkRepo := sqldb.NewShortLinkSql(sqlDB)
					gotIsExist, err := shortLinkRepo.IsAliasExist(testCase.alias)
					assert.Equal(t, nil, err)
					assert.Equal(t, testCase.expIsExist, gotIsExist)
				})
		})
	}
}

func TestShortLinkSql_GetShortLinkByAlias(t *testing.T) {
	twoYearsAgo := mustParseTime(t, "2017-05-01T08:02:16-07:00")
	now := mustParseTime(t, "2019-05-01T08:02:16-07:00")

	testCases := []struct {
		name              string
		tableRows         []shortLinkTableRow
		alias             string
		hasErr            bool
		expectedShortLink entity.ShortLink
	}{
		{
			name:      "alias not found",
			tableRows: []shortLinkTableRow{},
			alias:     "220uFicCJj",
			hasErr:    true,
		},
		{
			name: "found short link",
			tableRows: []shortLinkTableRow{
				{
					alias:     "220uFicCJj",
					longLink:  "http://www.google.com",
					createdAt: &twoYearsAgo,
					expireAt:  &now,
					updatedAt: &now,
				},
				{
					alias:     "yDOBcj5HIPbUAsw",
					longLink:  "http://www.facebook.com",
					createdAt: &twoYearsAgo,
					expireAt:  &now,
					updatedAt: &now,
				},
			},
			alias:  "220uFicCJj",
			hasErr: false,
			expectedShortLink: entity.ShortLink{
				Alias:     "220uFicCJj",
				LongLink:  "http://www.google.com",
				CreatedAt: &twoYearsAgo,
				ExpireAt:  &now,
				UpdatedAt: &now,
			},
		},
		{
			name: "nil time",
			tableRows: []shortLinkTableRow{
				{
					alias:     "220uFicCJj",
					longLink:  "http://www.google.com",
					createdAt: nil,
					expireAt:  nil,
					updatedAt: nil,
				},
				{
					alias:     "yDOBcj5HIPbUAsw",
					longLink:  "http://www.facebook.com",
					createdAt: &twoYearsAgo,
					expireAt:  &now,
					updatedAt: &now,
				},
			},
			alias:  "220uFicCJj",
			hasErr: false,
			expectedShortLink: entity.ShortLink{
				Alias:     "220uFicCJj",
				LongLink:  "http://www.google.com",
				CreatedAt: nil,
				ExpireAt:  nil,
				UpdatedAt: nil,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dbtest.AccessTestDB(
				dbConnector,
				dbMigrationTool,
				dbMigrationRoot,
				dbConfig,
				func(sqlDB *sql.DB) {
					insertShortLinkTableRows(t, sqlDB, testCase.tableRows)

					shortLinkRepo := sqldb.NewShortLinkSql(sqlDB)
					shortLink, err := shortLinkRepo.GetShortLinkByAlias(testCase.alias)

					if testCase.hasErr {
						assert.NotEqual(t, nil, err)
						return
					}
					assert.Equal(t, nil, err)
					assert.Equal(t, testCase.expectedShortLink, shortLink)
				},
			)
		})
	}
}

func TestShortLinkSql_CreateShortLink(t *testing.T) {
	now := mustParseTime(t, "2019-05-01T08:02:16-07:00")

	testCases := []struct {
		name      string
		tableRows []shortLinkTableRow
		shortLink entity.ShortLink
		hasErr    bool
	}{
		{
			name: "alias exists",
			tableRows: []shortLinkTableRow{
				{
					alias:    "220uFicCJj",
					longLink: "http://www.facebook.com",
					expireAt: &now,
				},
			},
			shortLink: entity.ShortLink{
				Alias:    "220uFicCJj",
				LongLink: "http://www.google.com",
				ExpireAt: &now,
			},
			hasErr: true,
		},
		{
			name: "successfully create short link",
			tableRows: []shortLinkTableRow{
				{
					alias:    "abc",
					longLink: "http://www.google.com",
					expireAt: &now,
				},
			},
			shortLink: entity.ShortLink{
				Alias:    "220uFicCJj",
				LongLink: "http://www.google.com",
				ExpireAt: &now,
			},
			hasErr: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dbtest.AccessTestDB(
				dbConnector,
				dbMigrationTool,
				dbMigrationRoot,
				dbConfig,
				func(sqlDB *sql.DB) {
					insertShortLinkTableRows(t, sqlDB, testCase.tableRows)

					shortLinkRepo := sqldb.NewShortLinkSql(sqlDB)
					err := shortLinkRepo.CreateShortLink(testCase.shortLink)

					if testCase.hasErr {
						assert.NotEqual(t, nil, err)
						return
					}
					assert.Equal(t, nil, err)
				},
			)
		})
	}
}

func TestShortLinkSql_UpdateShortLink(t *testing.T) {
	createdAt := mustParseTime(t, "2017-05-01T08:02:16-07:00")
	now := time.Now()

	testCases := []struct {
		name              string
		oldAlias          string
		newShortLink      entity.ShortLink
		tableRows         []shortLinkTableRow
		hasErr            bool
		expectedShortLink entity.ShortLink
	}{
		{
			name:     "alias not found",
			oldAlias: "does_not_exist",
			tableRows: []shortLinkTableRow{
				{
					alias:     "220uFicCJj",
					longLink:  "https://www.google.com",
					createdAt: &createdAt,
				},
			},
			hasErr:            true,
			expectedShortLink: entity.ShortLink{},
		},
		{
			name:     "alias is taken",
			oldAlias: "220uFicCja",
			tableRows: []shortLinkTableRow{
				{
					alias:     "220uFicCJj",
					longLink:  "https://www.google.com",
					createdAt: &createdAt,
				},
				{
					alias:     "efpIZ4OS",
					longLink:  "https://gmail.com",
					createdAt: &createdAt,
				},
			},
			hasErr:            true,
			expectedShortLink: entity.ShortLink{},
		},
		{
			name:     "valid new alias",
			oldAlias: "220uFicCJj",
			newShortLink: entity.ShortLink{
				Alias:     "GxtKXM9V",
				LongLink:  "https://www.google.com",
				UpdatedAt: &now,
			},
			tableRows: []shortLinkTableRow{
				{
					alias:     "220uFicCJj",
					longLink:  "https://www.google.com",
					createdAt: &createdAt,
				},
			},
			hasErr: false,
			expectedShortLink: entity.ShortLink{
				Alias:     "GxtKXM9V",
				LongLink:  "https://www.google.com",
				UpdatedAt: &now,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dbtest.AccessTestDB(
				dbConnector,
				dbMigrationTool,
				dbMigrationRoot,
				dbConfig,
				func(sqlDB *sql.DB) {
					insertShortLinkTableRows(t, sqlDB, testCase.tableRows)
					expectedShortLink := testCase.expectedShortLink

					shortLinkRepo := sqldb.NewShortLinkSql(sqlDB)
					shortLink, err := shortLinkRepo.UpdateShortLink(
						testCase.oldAlias,
						testCase.newShortLink,
					)

					if testCase.hasErr {
						assert.NotEqual(t, nil, err)
						return
					}
					assert.Equal(t, nil, err)
					assert.Equal(t, expectedShortLink.Alias, shortLink.Alias)
					assert.Equal(t, expectedShortLink.LongLink, shortLink.LongLink)
					assert.Equal(t, expectedShortLink.ExpireAt, shortLink.ExpireAt)
					assert.Equal(t, expectedShortLink.UpdatedAt, shortLink.UpdatedAt)
				},
			)
		})
	}
}

func TestShortLinkSql_GetShortLinkByAliases(t *testing.T) {
	twoYearsAgo := mustParseTime(t, "2017-05-01T08:02:16-07:00")
	now := mustParseTime(t, "2019-05-01T08:02:16-07:00")

	testCases := []struct {
		name               string
		tableRows          []shortLinkTableRow
		aliases            []string
		hasErr             bool
		expectedShortLinks []entity.ShortLink
	}{
		{
			name:      "alias not found",
			tableRows: []shortLinkTableRow{},
			aliases:   []string{"220uFicCJj"},
			hasErr:    false,
		},
		{
			name: "found short link",
			tableRows: []shortLinkTableRow{
				{
					alias:     "220uFicCJj",
					longLink:  "http://www.google.com",
					createdAt: &twoYearsAgo,
					expireAt:  &now,
					updatedAt: &now,
				},
				{
					alias:     "yDOBcj5HIPbUAsw",
					longLink:  "http://www.facebook.com",
					createdAt: &twoYearsAgo,
					expireAt:  &now,
					updatedAt: &now,
				},
			},
			aliases: []string{"220uFicCJj", "yDOBcj5HIPbUAsw"},
			hasErr:  false,
			expectedShortLinks: []entity.ShortLink{
				{
					Alias:     "220uFicCJj",
					LongLink:  "http://www.google.com",
					CreatedAt: &twoYearsAgo,
					ExpireAt:  &now,
					UpdatedAt: &now,
				},
				{
					Alias:     "yDOBcj5HIPbUAsw",
					LongLink:  "http://www.facebook.com",
					CreatedAt: &twoYearsAgo,
					ExpireAt:  &now,
					UpdatedAt: &now,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dbtest.AccessTestDB(
				dbConnector,
				dbMigrationTool,
				dbMigrationRoot,
				dbConfig,
				func(sqlDB *sql.DB) {
					insertShortLinkTableRows(t, sqlDB, testCase.tableRows)

					shortLinkRepo := sqldb.NewShortLinkSql(sqlDB)
					shortLink, err := shortLinkRepo.GetShortLinksByAliases(testCase.aliases)

					if testCase.hasErr {
						assert.NotEqual(t, nil, err)
						return
					}
					assert.Equal(t, nil, err)
					assert.Equal(t, testCase.expectedShortLinks, shortLink)
				},
			)
		})
	}
}

func insertShortLinkTableRows(t *testing.T, sqlDB *sql.DB, tableRows []shortLinkTableRow) {
	for _, tableRow := range tableRows {
		_, err := sqlDB.Exec(
			insertShortLinkRowSQL,
			tableRow.alias,
			tableRow.longLink,
			tableRow.createdAt,
			tableRow.expireAt,
			tableRow.updatedAt,
		)
		assert.Equal(t, nil, err)
	}
}
