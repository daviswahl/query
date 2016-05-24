package query

import (
	"database/sql/driver"
	"gitolite.corp.redhat.com/it-eng/go/assert"
	"testing"
	"time"
)

var selectTests []interface{} = []interface{}{Type.String(), Type.Int(), Type.Time()}

func TestSelect(t *testing.T) {
	for _, testCase := range selectTests {
		s, _ := Select("foo", testCase)
		assert.PtrTypeEq(t, testCase, s.scanner)
	}

	_, err := Select("foo", "")
	if err == nil {
		t.Logf("Expected error for unsupported sql type, but no error was returned")
		t.Fail()
	}
}

var values = [][]driver.Value{
	{int64(1), "12", time.Now()},
	{int64(12), "bar", time.Now()},
	{int64(0), "", time.Time{}},
}

func TestScanners(t *testing.T) {
	q := New().Select("s", Type.Int()).
		Select("s", Type.String()).
		Select("s", Type.Time())
	f := MockRows{}
	f.Cols = values
	i := 0
	for f.Next() {
		q.Scanners()
		f.Scan(q.Scanners()...)
		assert.SqlNullEq(t, values[i][0], q.Scanners()[0])
		assert.SqlNullEq(t, values[i][1], q.Scanners()[1])
		assert.SqlNullEq(t, values[i][2], q.Scanners()[2])
		i = i + 1
	}
}

type TestCase struct {
	expression string
	args       []interface{}
}

func TestWhere(t *testing.T) {
	w, _ := Where("WHERE BATS IN (?,?)", "foo", "bar")
	exp := []string{"foo", "bar"}
	for i, v := range w.args {
		assert.Eq(t, "", exp[i], v)
	}
}

func TestQueryWhere(t *testing.T) {
	q := New().Select("*", []driver.Value{}, "foo").
		Where("foo = ?", "bar").
		Where("qux IN (?, ?)", 12, 14)
	exp := []driver.Value{"bar", 12, 14}
	for i, v := range q.Args() {
		assert.Eq(t, "", exp[i], v)
	}
}

// All methods should return a new copy of a given query
func TestCopy(t *testing.T) {
	ae := func(e, a interface{}) { assert.Eq(t, "", e, a) }
	q := New()
	q.selects = []SelectExpression{SelectExpression{expression: "foo", scanner: Type.String()}}
	q = q.Limit(10, 5)
	ae("LIMIT 10, 5", q.limit.ToSql())
	q2 := q.Select("bar", Type.String())

	assert.Eq(t, "", 1, len(q.selects))
	assert.Eq(t, "", "LIMIT 10, 5", q.limit.ToSql())

	assert.Eq(t, "", 2, len(q2.selects))
	assert.Eq(t, "", "LIMIT 10, 5", q2.limit.ToSql())
}
