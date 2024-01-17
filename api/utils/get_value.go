package utils

import "database/sql"

type GetValue struct{}

func (gv GetValue) String(value *string) sql.NullString {
	if value != nil {
		return sql.NullString{
			String: *value,
			Valid:  true,
		}
	}

	return sql.NullString{}
}

func (gv GetValue) Int32(value *int32) sql.NullInt32 {
	if value != nil {
		return sql.NullInt32{
			Int32: *value,
			Valid:  true,
		}
	}

	return sql.NullInt32{}
}
