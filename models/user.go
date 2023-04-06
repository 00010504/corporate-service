package models

import "database/sql"

type NullShortUser struct {
	FirstName sql.NullString
	LastName  sql.NullString
	ID        sql.NullString
}
