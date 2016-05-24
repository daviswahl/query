package query

import (
	"fmt"
	//	"reflect"
)

type limit struct {
	Count  int
	Offset *int
}

func (l limit) ToSql() (str string) {
	if l.Count > 0 && l.Offset != nil {
		return fmt.Sprintf("LIMIT %v, %v", l.Count, *l.Offset)
	} else if l.Count > 0 {
		return fmt.Sprintf("LIMIT %v", l.Count)
	}
	return ""
}
