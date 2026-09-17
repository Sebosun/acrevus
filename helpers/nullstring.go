package helpers

import "database/sql"

func NewNullString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  true,
	}
}
