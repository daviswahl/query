package query

import (
	"database/sql"
	mysql "github.com/go-sql-driver/mysql"
	"gitolite.corp.redhat.com/it-eng/go/assert"
	"testing"
)

func TestString(t *testing.T) {
	assert.TypeEq(t, Type.String(), sql.NullString{})
}

func TestTime(t *testing.T) {
	assert.TypeEq(t, Type.Time(), mysql.NullTime{})
}

func TestInt(t *testing.T) {
	assert.TypeEq(t, Type.Int(), sql.NullInt64{})
}
