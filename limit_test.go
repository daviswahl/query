package query

import (
	"gitolite.corp.redhat.com/it-eng/go/assert"
	"testing"
)

func TestLimitOffset(t *testing.T) {
	i := 5
	e := "LIMIT 10, 5"
	a := limit{10, &i}.ToSql()
	assert.Eq(t, "", e, a)
}

func TestLimit(t *testing.T) {
	e := "LIMIT 10"
	a := limit{10, nil}.ToSql()
	assert.Eq(t, "", e, a)
}
