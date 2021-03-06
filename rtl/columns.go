package rtl

import (
	"strings"
)

type String struct {
	columnID int
	key      bool
	val      func(idx int) string
}

func StringColumn(id int, val func(int) string) String {
	return String{columnID: id, val: val}
}

func StringIndex(id int, val func(int) string) String {
	return String{columnID: id, key: true, val: val}
}

func (c String) query(val string, op test) Query {
	return queryForClause(clause{
		columnID: c.columnID,
		op:       op,
		cmp: func(idx int) int {
			return strings.Compare(c.val(idx), val)
		},
	})
}

func (c String) Eq(val string) Query {
	if c.key {
		return c.query(val, key)
	}
	return c.query(val, eq)
}

func (c String) Lt(val string) Query { return c.query(val, lt) }
func (c String) Le(val string) Query { return c.query(val, le) }
func (c String) Gt(val string) Query { return c.query(val, gt) }
func (c String) Ge(val string) Query { return c.query(val, ge) }
func (c String) Ne(val string) Query { return c.query(val, ne) }

func (c String) Range(from, to string) Query {
	return c.Ge(from).And(c.Le(to))
}

type Int struct {
	columnID int
	key      bool
	val      func(idx int) int
}

func IntColumn(id int, val func(int) int) Int {
	return Int{columnID: id, val: val}
}

func IntIndex(id int, val func(int) int) Int {
	return Int{columnID: id, key: true, val: val}
}

func (c Int) query(val int, op test) Query {
	return queryForClause(clause{
		columnID: c.columnID,
		op:       op,
		cmp: func(idx int) int {
			return c.val(idx) - val
		},
	})
}

func (c Int) Eq(val int) Query {
	if c.key {
		return c.query(val, key)
	}
	return c.query(val, eq)
}

func (c Int) Lt(val int) Query { return c.query(val, lt) }
func (c Int) Le(val int) Query { return c.query(val, le) }
func (c Int) Gt(val int) Query { return c.query(val, gt) }
func (c Int) Ge(val int) Query { return c.query(val, ge) }
func (c Int) Ne(val int) Query { return c.query(val, ne) }

func (c Int) Range(from, to int) Query {
	return c.Ge(from).And(c.Le(to))
}
