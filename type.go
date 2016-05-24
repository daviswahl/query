package query

import (
	"database/sql"
	mysql "github.com/go-sql-driver/mysql"
)

type (
	types int
)

const Type types = 0

func (t types) String() sql.NullString {
	return sql.NullString{}
}
func (t types) Time() mysql.NullTime {
	return mysql.NullTime{}
}

func (t types) Int() sql.NullInt64 {
	return sql.NullInt64{}
}
