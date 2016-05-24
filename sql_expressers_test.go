package query

import (
	"fmt"
	"gitolite.corp.redhat.com/it-eng/go/assert"
	"strings"
	"testing"
)

func TestLimitMethod(t *testing.T) {
	q := New()
	q = q.Limit(10, 5)
	assert.Eq(t, "", "LIMIT 10, 5", q.limit.ToSql())
}

func TestOrderByToSql(t *testing.T) {
	actual := orderBy("bar").ToSql()
	expected := "ORDER BY bar"
	assert.Eq(t, "", expected, actual)
}

func TestHavingToSql(t *testing.T) {
	actual := having("bar").ToSql()
	expected := "HAVING bar"
	assert.Eq(t, "", expected, actual)

}

func TestGroupByToSql(t *testing.T) {
	actual := groupBy("bar").ToSql()
	expected := "GROUP BY bar"
	assert.Eq(t, "", expected, actual)
}

func TestJoinToSql(t *testing.T) {
	actual := join("LEFT JOIN t1").ToSql()
	expected := "LEFT JOIN t1"
	assert.Eq(t, "", expected, actual)
}

func TestJoinsToSql(t *testing.T) {
	j1 := join("LEFT JOIN t1 ON t2.id = t1.id")
	j2 := join("BEST JOIN t2 ON t1.id = t2.id")
	expected := "\tLEFT JOIN t1 ON t2.id = t1.id\n\tBEST JOIN t2 ON t1.id = t2.id"
	assert.Eq(t, "", expected, joins{j1, j2}.ToSql())
}

func TestFromToSql(t *testing.T) {
	actual := from("bar").ToSql()
	expected := "FROM bar"
	assert.Eq(t, "", expected, actual)

}

func TestSelectOptionToSql(t *testing.T) {
	actual := selectOption("DISTINCT").ToSql()
	expected := "SELECT DISTINCT"
	assert.Eq(t, "", expected, actual)
}

func TestSelectExprToSql(t *testing.T) {
	s := SelectExpression{expression: "bar", scanner: Type.String()}
	actual := s.ToSql()
	expected := "bar"
	assert.Eq(t, "", expected, actual)
}

func TestSelectsToSql(t *testing.T) {
	e := "\tt1.bar,\n\tCOUNT(DISTINCT(t2.batz))"
	q := New()
	q = q.Select("t1.bar", Type.String())
	q = q.Select("COUNT(DISTINCT(t2.batz))", Type.Int())
	if a := q.selects.ToSql(); e != a {
		t.Logf("\nExpected:\n%v\nActual:\n%v", e, a)
		t.Fail()
	}
}

func TestQueryToSql(t *testing.T) {
	q := New().
		Select("t1.bar", Type.String()).
		Select("t2.batz", Type.Int())
	q = q.Limit(10, 5)
	q = q.From("t1")
	q = q.Join("LEFT JOIN t2 ON t2.t1_id = t1.id")
	q = q.GroupBy("t1.uuid")
	exprs := []string{
		q.selects.ToSql(),
		q.from.ToSql(),
		q.joins.ToSql(),
		q.groupBy.ToSql(),
		q.limit.ToSql(),
	}
	expected := fmt.Sprintf("SELECT\n%v", strings.Join(exprs, "\n"))
	expected = strings.Trim(expected, "\n")
	assert.Eq(t, "", expected, strings.Trim(q.ToSql(), "\n"))
}
