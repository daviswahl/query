package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

type (
	orderBy      string
	having       string
	groupBy      string
	join         string
	from         string
	selectOption string

	joins   []join
	wheres  []where
	selects []SelectExpression
)

type SelectExpression struct {
	expression string
	scanner    driver.Value
}

type Query struct {
	selectOption selectOption
	selects      selects
	from         from
	joins        joins
	wheres       wheres
	groupBy      groupBy
	having       having
	orderBy      orderBy
	limit        limit

	ValueMap map[interface{}]*driver.Value
}

type SqlExpresser interface {
	ToSql() string
}

func New() Query {
	q := Query{}
	q.ValueMap = make(map[interface{}]*driver.Value, 0)
	return q
}

func (q Query) Limit(lim int, offset ...int) Query {
	if len(offset) == 0 || len(offset) > 1 {
		o := 0
		q.limit = limit{lim, &o}
	} else {
		q.limit = limit{lim, &offset[0]}
	}
	return q
}

func (q Query) From(s string) Query {
	q.from = from(s)
	return q
}
func (q Query) LeftJoin(s string) Query {
	return q.Join("LEFT JOIN " + s)
}

func (q Query) Join(s string) Query {
	q.joins = append(q.joins, join(s))
	return q
}

func (q Query) GroupBy(s string) Query {
	q.groupBy = groupBy(s)
	return q
}
func Select(s string, T driver.Value) (SelectExpression, error) {
	switch (T).(type) {
	case sql.NullString:
		return SelectExpression{s, new(sql.NullString)}, nil
	case sql.NullInt64:
		return SelectExpression{s, new(sql.NullInt64)}, nil
	case mysql.NullTime:
		return SelectExpression{s, new(mysql.NullTime)}, nil
	default:
		return SelectExpression{}, errors.New(fmt.Sprintf("Unrecognized sql value %v", T))
	}
}

func (q Query) Select(s string, T driver.Value, key ...interface{}) Query {
	sel, err := Select(s, T)
	if err != nil {
		return Query{}
	}
	selects := append(q.selects, sel)
	q.selects = selects

	if key != nil {
		q.ValueMap[key[0]] = &sel.scanner
	}
	return q
}

type where struct {
	expression string
	args       []interface{}
}

func Where(s string, T ...interface{}) (where, error) {
	return where{s, T}, nil
}

func (q Query) Args() []driver.Value {
	//q.wheres
	return nil
}
func (q Query) Where(s string, args ...driver.Value) Query {
	return q
}

func (q Query) Scanners() []driver.Value {
	v := make([]driver.Value, len(q.selects))
	for i, s := range q.selects {
		v[i] = s.scanner
	}
	return v
}

func (q Query) Values() []driver.Value {
	vals := make([]driver.Value, len(q.Scanners()))
	for i, s := range q.Scanners() {
		switch t := s.(type) {
		case *sql.NullString:
			vals[i] = t.String
		case *sql.NullInt64:
			vals[i] = t.Int64
		case *mysql.NullTime:
			vals[i] = t.Time
		}
	}
	return vals
}
