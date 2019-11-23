package db

import (
	"database/sql"
	"fmt"
	"short/app/adapter/db/table"
	"short/app/entity"
	"short/app/usecase/repository"
)

var _ repository.UserURLRelation = (*UserURLRelationSQL)(nil)

// UserURLRelationSQL accesses UserURLRelation information in user_url_relation
// table.
type UserURLRelationSQL struct {
	db *sql.DB
}

// CreateRelation establishes bi-directional relationship between a user and a
// url in user_url_relation table.
func (u UserURLRelationSQL) CreateRelation(user entity.User, url entity.URL) error {
	statement := fmt.Sprintf(`
INSERT INTO "%s" ("%s","%s")
VALUES ($1,$2)
`,
		table.UserURLRelation.TableName,
		table.UserURLRelation.ColumnUserEmail,
		table.UserURLRelation.ColumnURLAlias,
	)

	_, err := u.db.Exec(statement, user.Email, url.Alias)
	return err
}

// NewUserURLRelationSQL creates UserURLRelationSQL
func NewUserURLRelationSQL(db *sql.DB) UserURLRelationSQL {
	return UserURLRelationSQL{
		db: db,
	}
}
