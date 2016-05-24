package query

import (
	"fmt"
	"strings"
)

func (s orderBy) ToSql() string {
	if s != "" {
		return fmt.Sprintf("ORDER BY %v", s)
	}
	return ""
}

func (s having) ToSql() string {
	if s != "" {
		return fmt.Sprintf("HAVING %v", s)
	}
	return ""
}

func (s groupBy) ToSql() string {
	if s != "" {
		return fmt.Sprintf("GROUP BY %v", s)
	}
	return ""
}

func (s where) ToSql() string {
	if s.expression != "" {
		return fmt.Sprintf("WHERE %v", s)
	}
	return ""
}

func (s join) ToSql() string {
	if s != "" {
		return string(s)
	}
	return ""
}

func (s joins) ToSql() string {
	str := make([]string, len(s))
	for i, expr := range s {
		str[i] = expr.ToSql()
	}
	return "\t" + strings.Join(str, "\n\t")
}

func (s from) ToSql() string {
	if s != "" {
		return fmt.Sprintf("FROM %v", s)
	}
	return "FROM"
}

func (s selectOption) ToSql() string {
	if s != "" {
		return fmt.Sprintf("SELECT %v", s)
	}
	return "SELECT"
}

func (exp SelectExpression) ToSql() string {
	return string(exp.expression)
}

func (exp selects) ToSql() string {
	str := ""
	ss := make([]string, len(exp))
	for i, s := range exp {
		ss[i] = s.ToSql()
	}
	str = strings.Join(ss, ",\n\t")
	return "\t" + str
}

func (q Query) ToSql() string {
	fields := []SqlExpresser{q.selectOption, q.selects,
		q.from, q.joins, q.where, q.groupBy, q.having, q.orderBy, q.limit}

	sql := ""
	for _, exp := range fields {
		str := exp.ToSql()
		if str != "" {
			sql = sql + str + "\n"
		}
	}
	return sql
}
