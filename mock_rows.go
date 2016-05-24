package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	mysql "github.com/go-sql-driver/mysql"
	"time"
)

type MockRows struct {
	sql.Rows
	Cols     [][]driver.Value
	lastcols []driver.Value
}

func (rs *MockRows) Next() bool {
	if len(rs.Cols) > 0 {
		rs.lastcols = rs.Cols[0]
		rs.Cols = rs.Cols[1:]
		return true
	}
	return false
}

func (rs *MockRows) Scan(dest ...driver.Value) error {
	for i, sv := range rs.lastcols {
		err := convertAssign(dest[i], sv)
		if err != nil {
			return fmt.Errorf("sql: Scan error on column index %d: %v", i, err)
		}
	}
	return nil
}

var errNilPtr = errors.New("NilPionter")

func convertAssign(dest, src interface{}) error {
	// Common cases, without reflect.
	switch s := src.(type) {
	case string:
		switch d := dest.(type) {
		case *sql.NullString:
			if s == "" {
				d.Valid = false
				d.String = ""
			} else {
				d.String = s
				d.Valid = true
			}
			return nil
		}
	case int64:
		switch d := dest.(type) {
		case *sql.NullInt64:
			if s == 0 {
				d.Int64 = int64(0)
				d.Valid = false
			} else {
				d.Int64 = int64(s)
				d.Valid = true
			}
			return nil
		}
	case time.Time:
		switch d := dest.(type) {
		case *mysql.NullTime:
			if s.IsZero() {
				d.Time = time.Time{}
				d.Valid = false
			} else {
				d.Time = s
				d.Valid = true
			}
			return nil
		}
	}
	return nil
}
