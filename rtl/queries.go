package rtl

import (
	"sort"
	"strings"
)

type Query struct {
	clauses []clause
}

type QueryResult struct {
	this, end int
	q         Query
}

type test byte

const (
	eq test = iota
	lt
	le
	gt
	ge
	ne
	key
)

type clause struct {
	columnID int
	op       test
	cmp      func(int) int
}

// Query construction

func (c StringColumn) query(val string, op test) Query {
	return queryForClause(clause{
		columnID: c.ColumnID,
		op:       op,
		cmp: func(idx int) int {
			return strings.Compare(c.Val(idx), val)
		},
	})
}

func (c StringColumn) Eq(val string) Query {
	if c.Key {
		return c.query(val, key)
	}
	return c.query(val, eq)
}

func (c StringColumn) Lt(val string) Query { return c.query(val, lt) }
func (c StringColumn) Le(val string) Query { return c.query(val, le) }
func (c StringColumn) Gt(val string) Query { return c.query(val, gt) }
func (c StringColumn) Ge(val string) Query { return c.query(val, ge) }
func (c StringColumn) Ne(val string) Query { return c.query(val, ne) }

func (c StringColumn) Range(from, to string) Query {
	return c.Ge(from).And(c.Le(to))
}

func (c IntColumn) query(val int, op test) Query {
	return queryForClause(clause{
		columnID: c.ColumnID,
		op:       op,
		cmp: func(idx int) int {
			return c.Val(idx) - val
		},
	})
}

func (c IntColumn) Eq(val int) Query {
	if c.Key {
		return c.query(val, key)
	}
	return c.query(val, eq)
}

func (c IntColumn) Lt(val int) Query { return c.query(val, lt) }
func (c IntColumn) Le(val int) Query { return c.query(val, le) }
func (c IntColumn) Gt(val int) Query { return c.query(val, gt) }
func (c IntColumn) Ge(val int) Query { return c.query(val, ge) }
func (c IntColumn) Ne(val int) Query { return c.query(val, ne) }

func (c IntColumn) Range(from, to int) Query {
	return c.Ge(from).And(c.Le(to))
}
func (c clause) key() bool { return c.op == key }

func queryForClause(c clause) Query {
	q := Query{}
	q.clauses = append(q.clauses, c)
	return q
}

// Query composition

func (q Query) And(r Query) Query {
	var res Query
	res.clauses = append(res.clauses, q.clauses...)
	res.clauses = append(res.clauses, r.clauses...)
	return res
}

// Query evaluation

func EvalQuery(q Query, n int) *QueryResult {
	res := &QueryResult{0, n, q}
	res.refineSearchSpace()
	res.this--
	return res
}

func (r *QueryResult) refineSearchSpace() {
	clauses := r.q.clauses
	sort.Slice(clauses, func(i, j int) bool {
		if clauses[i].key() != clauses[j].key() {
			return clauses[i].key()
		}
		return clauses[i].columnID < clauses[j].columnID
	})
	for i, c := range clauses {
		if c.op != key || i != c.columnID {
			break
		}
		r.this = sort.Search(r.end-r.this, func(idx int) bool {
			return c.cmp(idx+r.this) >= 0
		}) + r.this
		r.end = sort.Search(r.end-r.this, func(idx int) bool {
			return c.cmp(idx+r.this) > 0
		}) + r.this
	}
}

func (r *QueryResult) This() int {
	return r.this
}

func (r *QueryResult) Next() bool {
	for {
		r.this++
		if r.this >= r.end {
			return false
		}
		if r.q.matches(r.this) {
			return true
		}
	}
}

func (q Query) matches(idx int) bool {
	for _, c := range q.clauses {
		if !c.matches(idx) {
			return false
		}
	}
	return true
}

func (c clause) matches(idx int) bool {
	cmp := c.cmp(idx)
	switch c.op {
	case key, eq:
		return cmp == 0
	case lt:
		return cmp < 0
	case le:
		return cmp <= 0
	case gt:
		return cmp > 0
	case ge:
		return cmp >= 0
	case ne:
		return cmp != 0
	}
	panic("invalid op")
}
